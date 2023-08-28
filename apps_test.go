package github

import (
	"testing"
)

var testPrivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEA8XMXJdGRAKZm5nNN7Y5Hsd68otnxr5apyiSQrja9Jkjq0exL
LohCS1EIk0/+pB/vQxz1kR0E7HjhJTrChZQmLGgb8nSE987UEinRljTF0yjZSnjE
keNV8emtSGh2bLKPwd21b9wSAGANoFouEL0K6lLU8kzBRJuHNF9Uu0MiFfFkJVm3
2g1tMeUmJM0NRiBjQ77wdCNugEU5A41yyr++COgfAoFy5sPAtcp45ssI4lg0mMxI
vsP+/6caO8GYZzQl9aOHiGJONySKrd3SW5W+xnhDTkZN1tCqOIkvFl5Lpg+Cl0UE
UZjuqb626xfixOWtMD/dvrWxaM8yjaZEg4NZjwIDAQABAoIBAQCDwc1I6vJYy3Vt
nlBRKQpAqw5/Q7VanznqQEffeamAYdwaT/q62spqdT7bvJR1laOoGP58gLx2GoSq
H0WVRUILi4hsp18EJ46clstzTzsAvtLMi9igz9DPoTfZQoAVUt+V6FmhQBNmtwPY
lD19DtwNAMSJsI7q1IBUeQ0w3zKTi8ZC5UAT/Ng3vznYX4+1PNd5RRrngSVwZ8KO
q770P8RIRUmVI8KXNGeL/MN9l3HvwA1La2m/43RLczqXv5MP5u079FOYkK40TxlN
y4TJ40Ja6vG8G+B5CC+Kgr3ZVtOC1HlCknBqLjemUTcL6rmDFQTkPJ9BW/+xN1fC
PhJpV5e5AoGBAPr1uD5QLSxtk6XneUzqh+Jhv6xoq2S+DfSKkR4/KU0xVcDE90Jr
oQY6wJLsHblQ0Pp2H6/06c1d7rpta2Qo30w46D2vL4g6rH6tykHFkBypdi/gt3+Z
VRHA+IvHXlMwWXsHTYr/UVJon9SyiLnzPmyoV1/AeX6mp73b7eD37dIVAoGBAPZM
eYva1Dub1AUdFWNlLTmbgX6NRqdi37j2cFFd9osgeWw3pWgsHa0vWjeuRxDh5NIZ
9vg9zMrXEu7KvCVYAhQxNqvJxwsoBHkHNRF9FustGQGW2W0TLhiuLijUodVqYTJW
9Bhmz+sxC1FPBChiOXO7X0lm1Fihj/sHVYATNDoTAoGBAOImYDu3IJ4yuKT+rO7F
QmKc149UW29TXVwLKq7pGBz54l7uoCr4tojYlQVRRY/j5g5uOCvmNnLcvO6+/9Go
i2EyvwYnQlwvE5asoeEXWcCabWjDxlh0Ipb3IINFzBiHL3uQny4s2mm64p1XraJ0
MsLUCLi+yD17jRmogPsEMQnpAoGBAL0bZKuT/iYydCzk8roZQgscMeYH9PqqONpc
JUrkGVsjOPd1FkQQs0x4sg1Ue24j8zu6Ad0CHk6Tqg68jI8jrpzwWGi4CWKwfBat
CPr/j2xMeQm2WASemGMMwZZKBGPHRQ+QoeRmdDfBtU3dnHShTjlk4TmLgXj3u4Pj
Uqt+kzgPAoGBAOeI2EoJSOtpHlJgXT5v8qvqCFdu/oiiS/i9d25CkL0AIlT0YJZu
VAeKVKieSln/vQZfuklfdcmREwcn7LiMmU7KeLm5ehsfUAtjT/9c4KOj6+/unrQZ
NAYPACY/P+sO1+RN3UkezoYdbgnperNKSMreQtrxL/0wkDaKdmUWm0ty
-----END RSA PRIVATE KEY-----
	`)

func TestNewGitHubApp(t *testing.T) {

	// Test Case 1: Empty configuration
	cfg := &GitHubAppConfig{}
	_, err := NewGitHubApp(cfg)
	if err == nil {
		t.Errorf("Expected error for empty configuration")
	}

	// Test Case 2: Missing Application ID
	cfg = &GitHubAppConfig{
		InstallationID: 123,
	}
	_, err = NewGitHubApp(cfg)
	if err == nil {
		t.Errorf("Expected error for missing Application ID")
	}

	// Test Case 3: Missing Installation ID
	cfg = &GitHubAppConfig{
		ApplicationID: 298674,
	}
	_, err = NewGitHubApp(cfg)
	if err == nil {
		t.Errorf("Expected error for missing Installation ID")
	}

	// Test Case 4: Missing Private Key
	cfg = &GitHubAppConfig{
		ApplicationID:  298674,
		InstallationID: 123,
	}
	_, err = NewGitHubApp(cfg)
	if err == nil {
		t.Errorf("Expected error for missing Private Key")
	}
}
