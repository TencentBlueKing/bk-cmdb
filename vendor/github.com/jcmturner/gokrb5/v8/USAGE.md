## Version 8 Usage

### Configuration
The gokrb5 libraries use the same krb5.conf configuration file format as MIT Kerberos, 
described [here](https://web.mit.edu/kerberos/krb5-latest/doc/admin/conf_files/krb5_conf.html).
Config instances can be created by loading from a file path or by passing a string, io.Reader or bufio.Scanner to the 
relevant method:
```go
import "github.com/jcmturner/gokrb5/v8/config"
cfg, err := config.Load("/path/to/config/file")
cfg, err := config.NewConfigFromString(krb5Str) //String must have appropriate newline separations
cfg, err := config.NewConfigFromReader(reader)
cfg, err := config.NewConfigFromScanner(scanner)
```
### Keytab files
Standard keytab files can be read from a file or from a slice of bytes:
```go
import 	"github.com/jcmturner/gokrb5/v8/keytab"
ktFromFile, err := keytab.Load("/path/to/file.keytab")
ktFromBytes, err := keytab.Parse(b)

```

---

### Kerberos Client
**Create** a client instance with either a password or a keytab.
A configuration must also be passed. Additionally optional additional settings can be provided.
```go
import 	"github.com/jcmturner/gokrb5/v8/client"
cl := client.NewWithPassword("username", "REALM.COM", "password", cfg)
cl := client.NewWithKeytab("username", "REALM.COM", kt, cfg)
```
Optional settings are provided using the functions defined in the ``client/settings.go`` source file.

**Login**:
```go
err := cl.Login()
```
Kerberos Ticket Granting Tickets (TGT) will be automatically renewed unless the client was created from a CCache.

A client can be **destroyed** with the following method:
```go
cl.Destroy()
```

#### Active Directory KDC and FAST negotiation
Active Directory does not commonly support FAST negotiation so you will need to disable this on the client.
If this is the case you will see this error:
```KDC did not respond appropriately to FAST negotiation```
To resolve this disable PA-FX-Fast on the client before performing Login().
This is done with one of the optional client settings as shown below:
```go
cl := client.NewWithPassword("username", "REALM.COM", "password", cfg, client.DisablePAFXFAST(true))
```

#### Authenticate to a Service

##### HTTP SPNEGO
Create the HTTP request object and then create an SPNEGO client and use this to process the request with methods that 
are the same as on a HTTP client.
If nil is passed as the HTTP client when creating the SPNEGO client the http.DefaultClient is used.
When creating the SPNEGO client pass the Service Principal Name (SPN) or auto generate the SPN from the request 
object by passing a null string "".
```go
r, _ := http.NewRequest("GET", "http://host.test.gokrb5/index.html", nil)
spnegoCl := spnego.NewClient(cl, nil, "")
resp, err := spnegoCl.Do(r)
```

##### Generic Kerberos Client
To authenticate to a service a client will need to request a service ticket for a Service Principal Name (SPN) and form 
into an AP_REQ message along with an authenticator encrypted with the session key that was delivered from the KDC along 
with the service ticket.

The steps below outline how to do this.
* Get the service ticket and session key for the service the client is authenticating to.
The following method will use the client's cache either returning a valid cached ticket, renewing a cached ticket with 
the KDC or requesting a new ticket from the KDC.
Therefore the GetServiceTicket method can be continually used for the most efficient interaction with the KDC.
```go
tkt, key, err := cl.GetServiceTicket("HTTP/host.test.gokrb5")
```

The steps after this will be specific to the application protocol but it will likely involve a client/server 
Authentication Protocol exchange (AP exchange).
This will involve these steps:

* Generate a new Authenticator and generate a sequence number and subkey:
```go
auth, _ := types.NewAuthenticator(cl.Credentials.Realm, cl.Credentials.CName)
etype, _ := crypto.GetEtype(key.KeyType)
auth.GenerateSeqNumberAndSubKey(key.KeyType, etype.GetKeyByteSize())
```
* Set the checksum on the authenticator
The checksum is an application specific value. Set as follows:
```go
auth.Cksum = types.Checksum{
		CksumType: checksumIDint,
		Checksum:  checksumBytesSlice,
	}
```
* Create the AP_REQ:
```go
APReq, err := messages.NewAPReq(tkt, key, auth)
```

Now send the AP_REQ to the service. How this is done will be specific to the application use case.

#### Changing a Client Password
This feature uses the Microsoft Kerberos Password Change protocol (RFC 3244). 
This is implemented in Microsoft Active Directory and in MIT krb5kdc as of version 1.7.
Typically the kpasswd server listens on port 464.

Below is example code for how to use this feature:
```go
cfg, err := config.Load("/path/to/config/file")
if err != nil {
	panic(err.Error())
}
kt, err := keytab.Load("/path/to/file.keytab")
if err != nil {
	panic(err.Error())
}
cl := client.NewWithKeytab("username", "REALM.COM", kt)
cl.WithConfig(cfg)

ok, err := cl.ChangePasswd("newpassword")
if err != nil {
	panic(err.Error())
}
if !ok {
	panic("failed to change password")
}
```

The client kerberos config (krb5.conf) will need to have either the kpassd_server or admin_server defined in the 
relevant [realms] section. For example:
```
REALM.COM = {
  kdc = 127.0.0.1:88
  kpasswd_server = 127.0.0.1:464
  default_domain = realm.com
 }
```
See https://web.mit.edu/kerberos/krb5-latest/doc/admin/conf_files/krb5_conf.html#realms for more information.

#### Client Diagnostics
In the event of issues the configuration of a client can be investigated with its ``Diagnostics`` method.
This will check that the required enctypes defined in the client's krb5 config are available in its keytab.
It will also check that KDCs can be resolved for the client's REALM.
The error returned will contain details of any failed checks.
The configuration details of the client will be written to the ``io.Writer`` provided.

---

### Kerberised Service

#### SPNEGO/Kerberos HTTP Service
A HTTP handler wrapper can be used to implement Kerberos SPNEGO authentication for web services.
To configure the wrapper the keytab for the SPN and a Logger are required:
```go
kt, err := keytab.Load("/path/to/file.keytab")
l := log.New(os.Stderr, "GOKRB5 Service: ", log.Ldate|log.Ltime|log.Lshortfile)
```
Create a handler function of the application's handling method (apphandler in the example below):
```go
h := http.HandlerFunc(apphandler)
```
Configure the HTTP handler:
```go
http.Handler("/", spnego.SPNEGOKRB5Authenticate(h, &kt, service.Logger(l)))
```
The handler to be wrapped and the keytab are required arguments. 
Additional optional settings can be provided, such as the logger shown above.

Another example of optional settings may be that when using Active Directory where the SPN is mapped to a user account 
the keytab may contain an entry for this user account. In this case this should be specified as below with the 
``KeytabPrincipal``:
```go
http.Handler("/", spnego.SPNEGOKRB5Authenticate(h, &kt, service.Logger(l), service.KeytabPrincipal(pn)))
```

##### Session Management
For efficiency reasons it is not desirable to authenticate on every call to a web service. 
Therefore most authenticated web applications implement some form of session with the user.
Such sessions can be supported by passing a "session manager" into the ``SPNEGOKRB5Authenticate`` wrapper handler.
In order to not demand a specific session manager solution, the session manager must implement a simple interface:
```go
type SessionMgr interface {
	New(w http.ResponseWriter, r *http.Request, k string, v []byte) error
	Get(r *http.Request, k string) ([]byte, error)
}
```
- New - creates a new session for the request and adds a piece of data (key/value pair) to the session
- Get - extract from an existing session the value held within it under the key provided. 
This should return nil bytes or an error if there is no existing session.

The session manager (sm) that implements this interface should then be passed to the ``SPNEGOKRB5Authenticate`` wrapper 
handler as below:
```go
http.Handler("/", spnego.SPNEGOKRB5Authenticate(h, &kt, service.Logger(l), service.SessionManager(sm)))
```

The ``httpServer.go`` source file in the examples directory shows how this can be used with the popular gorilla web toolkit.

##### Validating Users and Accessing Users' Details
If authentication succeeds then the request's context will have a credentials objected added to it.
This object implements the ``github.com/jcmturner/goidentity/identity`` interface.
If Microsoft Active Directory is used as the KDC then additional ADCredentials are available in the 
``credentials.Attributes`` map under the key ``credentials.AttributeKeyADCredentials``. 
For example the SIDs of the users group membership are available and can be used by your application for authorization.

Checking and access the credentials within your application:
```go
// Get a goidentity credentials object from the request's context
creds := goidentity.FromHTTPRequestContext(r)
// Check if it indicates it is authenticated
if creds != nil && creds.Authenticated() {
    // Check for Active Directory attributes
	if ADCredsJSON, ok := creds.Attributes()[credentials.AttributeKeyADCredentials]; ok {
		ADCreds := new(credentials.ADCredentials)
        // Unmarshal the AD attributes
		err := json.Unmarshal([]byte(ADCredsJSON), ADCreds)
		if err == nil {
			// Now access the fields of the ADCredentials struct. For example: ADCreds.GroupMembershipSIDs
		}
	}
} else {
    // Not authenticated user
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprint(w, "Authentication failed")
}
```

#### Generic Kerberised Service - Validating Client Details
To validate the AP_REQ sent by the client on the service side call this method:
```go
import 	"github.com/jcmturner/gokrb5/v8/service"
s := service.NewSettings(&kt) // kt is a keytab and optional settings can also be provided.
if ok, creds, err := service.VerifyAPREQ(&APReq, s); ok {
        // Perform application specific actions
        // creds object has details about the client identity
}
```
