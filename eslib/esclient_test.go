package eslib

import "testing"

func TestSetFromUrl_Domain(t *testing.T) {
	c := NewEsClient()
	c.SetFromFlag("http://localhost")
	exp := "localhost"
	if c.Domain != exp {
		t.Errorf("Expected Domain: %s, but it was %s instead.", exp, c.Domain)
	}
}
func TestSetFromUrl_Port(t *testing.T) {
	c := NewEsClient()
	c.SetFromFlag("http://localhost:9000")
	exp := "9000"
	if c.Port != exp {
		t.Errorf("Expected Port: %s, but it was %s instead.", exp, c.Port)
	}
}
func TestSetFromUrl_Username(t *testing.T) {
	c := NewEsClient()
	c.SetFromFlag("http://user@localhost")
	exp := "user"
	if c.Username != exp {
		t.Errorf("Expected Username: %s, but it was %s instead.", exp, c.Username)
	}
}
func TestSetFromUrl_Password(t *testing.T) {
	c := NewEsClient()
	c.SetFromFlag("http://user:pass@localhost")
	exp := "pass"
	if c.Password != exp {
		t.Errorf("Expected Password: %s, but it was %s instead.", exp, c.Password)
	}
}
func TestSetFromUrl_EmptyURL(t *testing.T) {
	c := NewEsClient()
	err := c.SetFromFlag("")
	exp := "Url is empty"
	if err.Error() != exp {
		t.Errorf("Expected Error: %s, but it was %s instead.", exp, err)
	}
}

func TestNewRequest_basic(t *testing.T) {
	c := NewEsClient()
	c.SetFromFlag("http://localhost:9200")
	request, _ := c.NewRequest("GET", "_template", "")
	exp := "http://localhost:9200/_template"
	if request.URL.String() != exp {
		t.Errorf("Expected URL: %s, but it was %s instead.", exp, request.URL.String())
	}
}
func TestNewRequest_trimSlash(t *testing.T) {
	c := NewEsClient()
	c.SetFromFlag("http://localhost:9200")
	request, _ := c.NewRequest("GET", "/_template/", "")
	exp := "http://localhost:9200/_template"
	if request.URL.String() != exp {
		t.Errorf("Expected URL: %s, but it was %s instead.", exp, request.URL.String())
	}
}
func TestNewRequest_queryString(t *testing.T) {
	c := NewEsClient()
	c.SetFromFlag("http://localhost:9200")
	request, _ := c.NewRequest("GET", "_template", "wait=true")
	exp := "http://localhost:9200/_template?wait=true"
	if request.URL.String() != exp {
		t.Errorf("Expected URL: %s, but it was %s instead.", exp, request.URL.String())
	}
}
func TestNewRequest_basicAuth(t *testing.T) {
	c := NewEsClient()
	c.SetFromFlag("http://user:pass@localhost:9200")
	request, _ := c.NewRequest("GET", "_template", "")

	exp := "user"
	user, pass, valid := request.BasicAuth()
	if user != exp {
		t.Errorf("Expected Username: %s, but it was %s instead.", exp, user)
	}
	exp = "pass"
	if pass != exp {
		t.Errorf("Expected Password: %s, but it was %s instead.", exp, pass)
	}
	if !valid {
		t.Errorf("BasicAuth is not valid")
	}
}
func TestNewRequest_jsonHeader(t *testing.T) {
	c := NewEsClient()
	c.SetFromFlag("http://localhost:9200")
	request, _ := c.NewRequest("GET", "_template", "")
	exp := "application/json"
	if request.Header["Accept"][0] != exp {
		t.Errorf("Expected Header Accept: %s, but it was %s instead.", exp, request.Header["Accept"][0])
	}
}
