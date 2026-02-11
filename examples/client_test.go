package examples

import (
	"fmt"
	"testing"

	"ilicense-client-go/ilicense"
)

const publicKey = `MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsFlFkjtja1FFZVcfUBQrVjEABf3GUINizCGaCSQ9SHBvU+OUFUZbTujEckny+zejkErx5+KjVENbzd4/P3HkiH7X4dPiBRdDSFSrrLVG4UuCz7fLJ4BNwqnCa/lt4ZRgpYUsLWUWpnhXzWGQZgUYW+6zqqpiq+5hWNIeaRpDw3e8IJdXwFyEonUm9/52GwrlJSN4VAqYWQwEHxayF/eCGyBRzDU7AWhxSoWOoafLV3vnVOXr/h4myfHoQQYUu2oBVcLorxl4VPl2lXwcvRfxE2M8U++xTKSt7JC0t6j2Q28yWszjLubS3C4XUW3yLSvFPWdea0o7WmWnghjEagv0zQIDAQAB`
const code = "AAABU XsiaX NzdWV yX2Nv ZGUiO iJpc3 N1ZXI tYS1j b2RlI iwiaX NzdWV yX25h bWUiO iJpc3 N1ZXI tYS1u YW1lI iwiY3 VzdG9 tZXJf Y29kZ SI6Im N1c3R vbWVy LWEiL CJjdX N0b21 lcl9u YW1lI joiY3 VzdG9 tZXIt YS1uY W1lIi wicHJ vZHVj dF9jb 2RlIj oicHJ vZHVj dC1kI iwicH JvZHV jdF9u YW1lI joicH JvZHV jdC1k LW5hb WUiLC JsaWN lbnNl X2NvZ GUiOi JMSUM tMTc3 MDcxN zQxOD IyOCI sImlz c3VlX 2F0Ij oiMjA yNi0w Mi0xM FQxNz o1Njo 1OC4y MjgzN zErMD g6MDA iLCJl eHBpc mVfYX QiOiI yMDI2 LTA0L TEwVD E5OjA 1OjM1 WiIsI m1vZH VsZXM iOiJt LWEsb S1iIi wibWF 4X2lu c3Rhb mNlcy I6MH0 AAAEA rh3FM kUXme 0E-dE b0wbQ 4Lkr- qMIza vheKH RJI7q cJH6C xoubr 207Pu eT4M7 V5AEt PccRl fjmOk -efNL vduV4 G42QS B3pB7 n05nh yxfMH LcLpk It6wn 44Axo _Ms-f 5biD2 JnJwa py_AK 3RRa3 Hks7u FYaDZ UtXrB xxfyv YdUjR 55JYt gkzJW mokua DmpA3 XM6ZN LY7xk 0v4lA 1tIvp UR6zM 9ZK8R bCGrz hh5nk 9cExi y50rs xhaS5 f690e WIZYB XjZVN 3tkd4 5HZAT 3nSlS nzo-k W0HGS anDd9 wY0x2 GrYCR FRCxZ mlwaH tacBB w4Osz B8sZY WGTQI bykih 1O428 3w"

func TestClient(t *testing.T) {
	config := ilicense.DefaultConfig()
	config.PublicKey = publicKey
	client := ilicense.NewClient(&config)
	res, err := client.Activate(code)
	fmt.Printf("---%+v, %+v\n", res, err)
	client.CheckLicense()
	client.CheckLicenseStatus()
	client.CheckModule("m-a")
}
