package test 

import(
	"fmt"
	"testing"
	"github.com/begopher/property"
)

func Test_func_Evaluation_panic_when_rule_is_nil(t *testing.T){
	anyProperty := delegate[string]{t:t}	
	defer func(){
		got := recover()
		if got == nil {
			t.Fatal("passing nil to argument 'rule' must cause panic")
		}
		expected := "property.Evaluation: cannot be created from nil rule"
		if got != expected {
			t.Errorf("Expected message is (%v) got (%v)", expected, got)
		}
	}()
	property.Evaluation[string](nil, anyProperty)
}

func Test_func_Evaluation_panic_when_property_is_nil(t *testing.T){
	anyRule := rule[int]{}
	defer func(){
		got := recover()
		if got == nil {
			t.Fatal("passing nil to argument 'property' must cause panic")
		}
		expected := "property.Evaluation: cannot be created from nil property"
		if got != expected {
			t.Errorf("Expected message is (%v) got (%v)", expected, got)
		}
	}()
	property.Evaluation[int](anyRule, nil)
}

func Test_evaluation_Change_delegate_to_underlying_rule_without_mutation(t *testing.T){
	table := []string{ "go", "golang" }
	for _, expected := range table {
		rule := rule[string]{
			evaluate: func(got string) error {
				if got != expected {
					t.Errorf("expected value is (%v) got (%v)", expected, got)
				}
				return nil
			},
		}
		delegation := delegate[string]{
			change: func(string) error{ return nil },
		}
		evaluation := property.Evaluation[string](rule, delegation)
		evaluation.Change(expected)	
	}
}

func Test_evaluation_Change_rule_violation_prevents_delegation_to_underlying_property(t *testing.T){
	table := []string{ "any", "value" }
	for _, data := range table {
		rule := rule[string]{
			evaluate: func(string) error {
				return fmt.Errorf("Any error")
			},
		}
		delegation := delegate[string]{
			change: func(string) error {
				t.Errorf("delegation did occur to underlying property")
				return nil
			},
		}
		eval := property.Evaluation[string](rule, delegation)
		eval.Change(data)
	}
}

func Test_evaluation_Change_delegates_to_underlying_property_without_mutation(t *testing.T){
	table := []string{ "go", "golang" }
	for _, expected := range table {
		rule := rule[string]{
			evaluate: func(string) error {
				return nil
			},
		}
		delegation := delegate[string]{
			change: func(got string) error {
				if got != expected {
					t.Errorf("expected value is (%v) got (%v)", expected, got)
				}
				return nil
			},
		}
		eval := property.Evaluation[string](rule, delegation)
		eval.Change(expected)
	}
}

func Test_evaluation_Change_returns_error_from_underlying_property_without_mutation(t *testing.T){
	table := []error{
		nil,
		fmt.Errorf("Any custom error"),
	}
	for _, expected := range table {
		rule := rule[int]{
			evaluate: func(int) error {
				return nil
			},
		}
		delegation := delegate[int]{
			t:t,
			change: func(int) error {
				return expected
			},
		}
		eval := property.Evaluation[int](rule, delegation)
		var any int
		got := eval.Change(any)
		if got != expected {
			t.Errorf("expected error is (%v) got (%v)", expected, got)
		}
	}
}

func Test_evaluation_Value_delegates_to_underlying_property(t *testing.T){
	rule := rule[string]{
		evaluate: func(string) error {
			return nil
		},
	}
	var delegated bool
	delegation := delegate[string]{
		t:t,
		value: func() (string, error){
			delegated = true
			return "", nil
		},
	}
	eval := property.Evaluation[string](rule, delegation)
	eval.Value()
	if !delegated {
		t.Error("Delegation did not occur")
	}
}

func Test_evaluation_Value_returns_value_without_mutation(t *testing.T){
	table := []struct{
		value string
		err error
	}{
		{"Go", nil},
		{"any", fmt.Errorf("Any error")},
	}
	for _, data := range table {
		rule := rule[string]{
			evaluate: func(string) error {
				return nil
			},
		}
		delegation := delegate[string]{
			t:t,
			value: func() (string, error){
				return data.value, data.err
			},
		}
		eval := property.Evaluation[string](rule, delegation)
		got, _ := eval.Value()
		expected := data.value
		if got != expected {
			t.Errorf("expected value is (%v) got (%v)", expected, got)
		}
	}
}

func Test_evaluation_Value_returns_error_without_mutation(t *testing.T){
	table := []struct{
		value string
		err error
	}{
		{"Go", nil},
		{"any", fmt.Errorf("Any error")},
	}
	for _, data := range table {
		rule := rule[string]{
			evaluate: func(string) error {
				return nil
			},
		}
		delegation := delegate[string]{
			t:t,
			value: func() (string, error){
				return data.value, data.err
			},
		}
		eval := property.Evaluation[string](rule, delegation)
		_, got := eval.Value()
		expected := data.err
		if got != expected {
			t.Errorf("expected error is (%v) got (%v)", expected, got)
		}
	}
}
