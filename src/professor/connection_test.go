package professor

import (
	"testing"	
)

func TestSetKeyspace(t *testing.T) {

	c, err := Connect("127.0.0.1")
	defer c.Close()
	

	if err != nil {				
		t.Error(err)	
	} else {
		
		result, err := c.Query("use system")
		
		if err != nil {
			t.Error(err)			
		} else {
			if result != "system" {
				t.Errorf("Expected %s as a keyspace but got %s", "system", result)
			}
		}
	}
	
}