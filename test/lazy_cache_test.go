package test

import(
	"fmt"
	"testing"

	"github.com/begopher/property"
)

func Test_func_LazyCache_panic_when_property_is_nil(t *testing.T){
	defer func(){
		got := recover()
		if got == nil {
			t.Fatal("property.LazyCache: should panic when proeprty is nil")
		}
		expected := "property.LazyCache: cannot be created from nil property"
		if expected != got {
			t.Errorf("expected message is (%v) got (%v)", expected, got)
		}
	}()
	property.LazyCache[string](nil)
}

func Test_lazyCache_Change_delegates_to_underlying_property(t *testing.T){
	var invoked bool
	delegation := delegate[string]{
		t:t,
		change: func(string) error {
			invoked = true
			return nil
		},
	}
	cacheable := property.LazyCache[string](delegation)
	var any string
	cacheable.Change(any)
	if !invoked {
		t.Errorf("delegation did not occur to underlying property")
	}
}

func Test_cacheable_Change_does_not_mutate_value(t *testing.T){
	
}

func Test_lazyCache_Change_returns_error_of_underlying_property(t *testing.T){
	table := []error {
		nil,
		fmt.Errorf("any error"),
	}
	for _, expected := range table{
		delegation := delegate[string]{
			t:t,
			change: func(string) error {
				return expected
			},
		}
		cacheable := property.LazyCache[string](delegation)
		var any string
		got := cacheable.Change(any)
		if got != expected {
			t.Errorf("expected error is (%v) got (%v)", expected, got)
		}
	}
}


//
func Test_cacheable_Value_delegates_to_underlying_property_after_creation(t *testing.T){
	//cacheable.Value()
	//if !invoked {
		
	//}
}

