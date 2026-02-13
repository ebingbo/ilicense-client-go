package ilicense

import "fmt"

func ExampleClient_CheckLicenseStatus() {
	client := NewClient(nil)
	status, err := client.CheckLicenseStatus()
	fmt.Println(status, err == ErrLicenseNotFound)
	// Output: not_activated true
}
