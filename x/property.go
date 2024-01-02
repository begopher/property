package x

import(
	"github.com/begopher/rule"
	"github.com/begopher/rule/constraints"
	"github.com/begopher/event"
	"github.com/begopher/event/dispatcher"
)

type Property[T any] interface{
	Change(T) error
	Value() (T, error)
	Constraints(rule.Constraint[T])
	Cache(T)
	Publish(int, rule.Rule[T]) bool
	Unpublish(int) bool
	Bind(int, event.Registration) error
	Unbind(int, event.Registration) error
}


func Simple[T comparable](datasource Datasource[T]) Property[T] {
	if datasource == nil {
		panic("property.Simple: datasource cannot be nil")
	}
	return &property[T]{
		datasource: datasource,
		cons: constraints.New[T](),
		events: map[int]rule.Rule[T]{},
		dispatcher: dispatcher.New(),
	}
}
type property[T comparable] struct{
	// datasource to store and retrive value from
	datasource Datasource[T]
	// cons for entry validation
	cons       rule.Constraint[T]
	// rule for deciding when to publish an event
	events     map[int]rule.Rule[T]
	// dispatcher for sending notifications
	dispatcher event.Dispatcher
	cache *T
}

func (p *property[T]) Change(value T) error {
	if p.cache == nil {
		temp, err := p.Value()
		if err != nil {
			return err
		}
		p.cache = &temp
	}
	if *p.cache == value {
		return nil
	}
	if err := p.cons.Evaluate(value); err != nil {
		return err
	}
	if err := p.datasource.Change(value); err != nil {
		return err
	}
	for event, rule := range p.events {
		if rule.Evaluate(value) {
			p.dispatcher.Send(event)
		}
	}
	return nil
	
}

func (p *property[T]) Value() (T, error) {
	if p.cache != nil {
		return *p.cache, nil
	}
	temp, err := p.datasource.Value()
	if err == nil {
		p.cache =&temp
	}
	return temp, err
}

func (p *property[T]) Constraints(cons rule.Constraint[T]) {
	if cons == nil {
		p.cons = constraints.New[T]()
		return 
	}
	p.cons = cons
}

func (p *property[T]) Cache(value T) {
	p.cache =&value
}

func (p *property[T]) Publish(event int, r rule.Rule[T]) bool {
	if r == nil {
		return false
	}
	if !p.dispatcher.Publish(event) {
		return false
	}
	p.events[event] = r
	return true 
	
}

func (p *property[T]) Unpublish(event int) bool {
	if p.dispatcher.Unpublish(event) {
		delete(p.events, event)
		return true 
	}
	return false
}

func (p *property[T]) Bind(event int, reg event.Registration) error {
	return p.dispatcher.Bind(event, reg)
}

func (p *property[T]) Unbind(event int, reg event.Registration) error {
	return p.dispatcher.Unbind(event, reg)
}
