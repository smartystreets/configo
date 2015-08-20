# configo
--
    import "github.com/smartystreets/configo"

Package configo provides flexible configurtion for go applications from any
number of 'sources', each of which implement the Source interface also provided
by this package.

Sources provided in this package:

    - JSONSource (key/value pairs in JSON content)
    - EnvironmentSource (key/value pairs from the environment)
    - CommandLineSource (key/value pairs via command line flags)
    - DefaultSource (key/value pairs manually configured by the application)
    - ConditionalSource (filters key/value retreival based on a condition)

These basic sources have been composed into additional sources, made available
via the following constructor methods:

    - FromDevelopmentOnlyDefaults()
    - FromRequiredInProductionJSONFile()
    - NewDefaultCommandLineConfigFileSource()
    - NewCommandLineConfigFileSource(path string)

Any of these sources may be provided to a Reader which is then used to retrieve
configuration values based on keys contained in the sources.

The reader can fetch values of various types:

    func (*Reader) Strings(key string) []string
    func (*Reader) String(key string) string
    func (*Reader) Ints(key string) []int
    func (*Reader) Int(key string) int
    func (*Reader) Bool(key string) bool
    func (*Reader) URLs(key string) []net.url.URL
    func (*Reader) URL(key string) net.url.URL

For each of the types returned above there are different ways to handle the
scenario when a key is not found. I'll illustrate this with the applicable Int
functions (but similar methods are implemented for each returned type):

    // returns zero value if key not found or values are malformed.
    func (*Reader) Int(key string) int

    // returns the value or the specified default if the key is not found.
    func (*Reader) IntDefault(key string, Default int) int

    // returns 0 and an error if key not found or values are malformed.
    func (*Reader) IntError(key string) (int, error)

    // returns the value or panics if the key is not found or the values are malformed.
    func (*Reader) IntPanic(key string) int

    // returns the value or calls log.Fatal() if the key is not found or the values are malformed.
    func (*Reader) IntFatal(key string) int

Here's a full example:

    reader := configo.NewReader(
        NewDefaultCommandLineConfigFileSource(),
        NewCommandLineFlag("s3-storage-address", "The address of the s3 bucket"),
        FromOptionalJSONFile("config-prod.json"),
    )
    value := reader.URL("s3-storage-address")

## Usage

```go
var (
	KeyNotFoundError    = errors.New("The specified key was not found.")
	MalformedValueError = errors.New("The specified value could not be parsed.")
)
```

#### type CommandLineConfigFileSource

```go
type CommandLineConfigFileSource struct {
}
```

CommandLineConfigFileSource registers a command line flag for specifying an
optional json config file, beyond any other config file definitions that follow
this one. It is intened to be used to provide an orveride to the regularly used
config file(s), like when you might be debugging in production (admit it, you've
been there too).

#### func  NewCommandLineConfigFileSource

```go
func NewCommandLineConfigFileSource(flagName string) *CommandLineConfigFileSource
```
NewDefaultCommandLineConfigFileSource registers a command line flag with the
given flagName for specifying an alternate JSON config file.

#### func  NewDefaultCommandLineConfigFileSource

```go
func NewDefaultCommandLineConfigFileSource() *CommandLineConfigFileSource
```
NewDefaultCommandLineConfigFileSource registers a command line flag called
"config" for specifying an alternate JSON config file.

#### func (*CommandLineConfigFileSource) Initialize

```go
func (this *CommandLineConfigFileSource) Initialize()
```
Initialize parses the command line flag and reads the altnerate JSON source.

#### func (*CommandLineConfigFileSource) Strings

```go
func (this *CommandLineConfigFileSource) Strings(key string) ([]string, error)
```
Strings reads the key from the JSON source if it was successfully loaded during
Initialize.

#### type CommandLineSource

```go
type CommandLineSource struct {
}
```

CommandLineSource registers a single command line flag and stores it's actual
value, if supplied on the command line. It implements the Source interface so it
can be used by a Reader.

#### func  NewCommandLineFlag

```go
func NewCommandLineFlag(name string, description string) *CommandLineSource
```
NewCommandLineFlag receives the name, defaultValue, and description of a command
line flag. The default value can be of any type handled by the internal
convertString function.

#### func (*CommandLineSource) Initialize

```go
func (this *CommandLineSource) Initialize()
```
Initialize calls flag.Parse(). Do not call until all CommandLineSource instances
have been created.

#### func (*CommandLineSource) Strings

```go
func (this *CommandLineSource) Strings(key string) ([]string, error)
```
Strings returns the command line flag value, or the default if no value was
provided at the command line.

#### type ConditionalSource

```go
type ConditionalSource struct {
}
```

ConditionalSource resolves values based on a condition supplied as a callback.

#### func  FromDevelopmentOnlyDefaults

```go
func FromDevelopmentOnlyDefaults() *ConditionalSource
```
A conditional source which determines if we are running in a development
environment

#### func  NewConditionalSource

```go
func NewConditionalSource(condition func() bool) *ConditionalSource
```
NewConditionalSource creates a conditional source with the provided condition
callback.

#### func (*ConditionalSource) Add

```go
func (this *ConditionalSource) Add(key string, values ...interface{}) *ConditionalSource
```
Add registers a key/values pairing for retrieval, so long as the supplied
condition is true.

#### func (*ConditionalSource) Initialize

```go
func (this *ConditionalSource) Initialize()
```

#### func (*ConditionalSource) Strings

```go
func (this *ConditionalSource) Strings(key string) ([]string, error)
```
Strings returns the value of the corresponding key, or KeyNotFoundError if the
condition is false.

#### type DefaultSource

```go
type DefaultSource struct {
}
```

DefaultSource is allows registration of specified default values of various
types.

#### func  NewDefaultSource

```go
func NewDefaultSource() *DefaultSource
```
NewDefaultSource initializes a new DefaultSource.

#### func (*DefaultSource) Add

```go
func (this *DefaultSource) Add(key string, values ...interface{}) *DefaultSource
```
Adds the provided values (which will be converted to strings) to the given key.
It does NOT overwrite existing values, it adds.

#### func (*DefaultSource) Initialize

```go
func (this *DefaultSource) Initialize()
```

#### func (*DefaultSource) Strings

```go
func (this *DefaultSource) Strings(key string) ([]string, error)
```
Strings returns all values associated with the given key, or KeyNotFoundError.

#### type EnvironmentSource

```go
type EnvironmentSource struct {
}
```

EnvironmentSource reads key-value pairs from the environment.

#### func  FromEnvironment

```go
func FromEnvironment() *EnvironmentSource
```
FromEnvironment creates an envirnoment source capable of parsing values
separated by the pipe character.

#### func  FromEnvironmentCustomSeparator

```go
func FromEnvironmentCustomSeparator(prefix, separator string) *EnvironmentSource
```
FromEnvironmentWithPrefix creates an envirnoment source capable of parsing
values separated by the specified character.

#### func  FromEnvironmentWithPrefix

```go
func FromEnvironmentWithPrefix(prefix string) *EnvironmentSource
```
FromEnvironmentWithPrefix creates an envirnoment source capable of: - reading
values with keys all beginning with the provided prefix, - parsing values
separated by the pipe character.

#### func (*EnvironmentSource) Initialize

```go
func (this *EnvironmentSource) Initialize()
```

#### func (*EnvironmentSource) Strings

```go
func (this *EnvironmentSource) Strings(key string) ([]string, error)
```
Strings reads the environment variable specified by key and returns the value or
KeyNotFoundError.

#### type JSONSource

```go
type JSONSource struct {
}
```

JSONSource houses key-value pairs unmarshaled from JSON data.

#### func  FromConditionalJSONFile

```go
func FromConditionalJSONFile(filename string, condition func() bool) *JSONSource
```
If the provided condition returns true, the specified filename is required and
must be found; otherwise loading the file is optional.

#### func  FromJSONContent

```go
func FromJSONContent(raw []byte) *JSONSource
```
FromJSONContent unmarshals the provided json content into a JSONSource. Any
resulting error results in a panic.

#### func  FromJSONFile

```go
func FromJSONFile(filename string) *JSONSource
```
FromJSONFile reads and unmarshals the file at the provided path into a
JSONSource. Any resulting error results in a panic.

#### func  FromOptionalJSONFile

```go
func FromOptionalJSONFile(filename string) *JSONSource
```
FromOptionalJSONFile is like FromJSONFile but it does not panic if the file is
not found.

#### func (*JSONSource) Initialize

```go
func (this *JSONSource) Initialize()
```

#### func (*JSONSource) Strings

```go
func (this *JSONSource) Strings(key string) ([]string, error)
```

#### type NoopSource

```go
type NoopSource struct{}
```


#### func (NoopSource) Initialize

```go
func (this NoopSource) Initialize()
```

#### func (NoopSource) Strings

```go
func (this NoopSource) Strings(string) ([]string, error)
```

#### type Reader

```go
type Reader struct {
}
```

Reader retrieves values from the provided sources, handling conversions to the
type identified by the method being called (Strings, Ints, etc...).

#### func  NewReader

```go
func NewReader(sources ...Source) *Reader
```
NewReader initializes a new reader using the provided sources. It calls each
non-nil source's Initialize() method.

#### func (*Reader) Bool

```go
func (this *Reader) Bool(key string) bool
```
Bool returns the boolean value associated with the given key or false if the key
does not exist or the value could not be parsed as a bool.

#### func (*Reader) BoolDefault

```go
func (this *Reader) BoolDefault(key string, Default bool) bool
```
BoolDefault returns the boolean value associated with the given key or returns
the provided default if the key does not exist or the value could not be parsed
as a bool.

#### func (*Reader) BoolError

```go
func (this *Reader) BoolError(key string) (bool, error)
```
BoolError returns the boolean value associated with the given key with an error
if the key does not exist or the value could not be parsed as a bool (according
to strconv.ParseBool).

#### func (*Reader) BoolFatal

```go
func (this *Reader) BoolFatal(key string) bool
```
BoolFatal returns the boolean value associated with the given key or calls
log.Fatal() if the key does not exist or the value could not be parsed as a
bool.

#### func (*Reader) BoolPanic

```go
func (this *Reader) BoolPanic(key string) bool
```
BoolPanic returns the boolean value associated with the given key or panics if
the key does not exist or the value could not be parsed as a bool.

#### func (*Reader) Int

```go
func (this *Reader) Int(key string) int
```
Int returns the first integer value associated with the given key or returns 0
if the key does not exist.

#### func (*Reader) IntDefault

```go
func (this *Reader) IntDefault(key string, Default int) int
```
IntDefault returns the first integer values associated with the given key or
returns the provided default if the key does not exist or the values could not
be parsed as integers.

#### func (*Reader) IntError

```go
func (this *Reader) IntError(key string) (int, error)
```
IntError returns the first integer value associated with the given key with an
error if the key does not exist or the values could not be parsed as integers
(according to strconv.Atoi).

#### func (*Reader) IntFatal

```go
func (this *Reader) IntFatal(key string) int
```
IntFatal returns the first integer value associated with the given key or calls
log.Fatal() if the key does not exist or the values could not be parsed as
integers.

#### func (*Reader) IntPanic

```go
func (this *Reader) IntPanic(key string) int
```
IntPanic returns the first integer value associated with the given key or panics
if the key does not exist or the values could not be parsed as integers.

#### func (*Reader) Ints

```go
func (this *Reader) Ints(key string) []int
```
Ints returns all integer values associated with the given key or returns 0 if
the key does not exist.

#### func (*Reader) IntsDefault

```go
func (this *Reader) IntsDefault(key string, Default []int) []int
```
IntsDefault returns all integer values associated with the given key or returns
provided defaults if the key does not exist or the values could not be parsed as
integers.

#### func (*Reader) IntsError

```go
func (this *Reader) IntsError(key string) ([]int, error)
```
IntsError returns all integer values associated with the given key with an error
if the key does not exist or the values could not be parsed as integers.

#### func (*Reader) IntsFatal

```go
func (this *Reader) IntsFatal(key string) []int
```
IntsFatal returns all integer values associated with the given key or calls
log.Fatal() if the key does not exist or the values could not be parsed as
integers.

#### func (*Reader) IntsPanic

```go
func (this *Reader) IntsPanic(key string) []int
```
IntsPanic returns all integer values associated with the given key or panics if
the key does not exist or the values could not be parsed as integers.

#### func (*Reader) String

```go
func (this *Reader) String(key string) string
```
String returns the first value associated with the given key or an empty string
if the key does not exist.

#### func (*Reader) StringDefault

```go
func (this *Reader) StringDefault(key string, Default string) string
```
StringDefault returns the first value associated with the given key or the
provided default if the key does not exist.

#### func (*Reader) StringError

```go
func (this *Reader) StringError(key string) (string, error)
```
StringError returns the first value associated with the given key with an an
error if the key does not exist.

#### func (*Reader) StringFatal

```go
func (this *Reader) StringFatal(key string) string
```
StringFatal returns the first value associated with the given key or calls
log.Fatal() if the key does not exist.

#### func (*Reader) StringPanic

```go
func (this *Reader) StringPanic(key string) string
```
StringPanic returns the first value associated with the given key or panics if
the key does not exist.

#### func (*Reader) Strings

```go
func (this *Reader) Strings(key string) []string
```
Strings returns all values associated with the given key or nil if the key does
not exist.

#### func (*Reader) StringsDefault

```go
func (this *Reader) StringsDefault(key string, Default []string) []string
```
StringsDefault returns all values associated with the given key or the provided
defaults if the key does not exist.

#### func (*Reader) StringsError

```go
func (this *Reader) StringsError(key string) ([]string, error)
```
StringsError returns all values associated with the given key with an error if
the key does not exist. It does so by searching it sources, in the order they
were provided, and returns the first non-error result or KeyNotFoundError.

#### func (*Reader) StringsFatal

```go
func (this *Reader) StringsFatal(key string) []string
```
StringsFatal returns all values associated with the given key or log.Fatal() if
the key does not exist.

#### func (*Reader) StringsPanic

```go
func (this *Reader) StringsPanic(key string) []string
```
StringsPanic returns all values associated with the given key or panics if the
key does not exist.

#### func (*Reader) URL

```go
func (this *Reader) URL(key string) url.URL
```
URL returns the first URL associated with the given key or returns the zero
value if the key does not exist or the value could not be parsed as a URL.

#### func (*Reader) URLDefault

```go
func (this *Reader) URLDefault(key string, Default url.URL) url.URL
```
URLDefault returns the first URL associated with the given key or returns
provided defaults if the key does not exist or the values could not be parsed as
URLs.

#### func (*Reader) URLError

```go
func (this *Reader) URLError(key string) (url.URL, error)
```
URLError returns the first URL associated with the given key with an error if
the key does not exist or the values could not be parsed as URLs.

#### func (*Reader) URLFatal

```go
func (this *Reader) URLFatal(key string) url.URL
```
URLFatal returns the first URL associated with the given key or calls
log.Fatal() if the key does not exist or the values could not be parsed as URLs.

#### func (*Reader) URLPanic

```go
func (this *Reader) URLPanic(key string) url.URL
```
URLPanic returns the first URL associated with the given key or panics if the
key does not exist or the values could not be parsed as URLs.

#### func (*Reader) URLs

```go
func (this *Reader) URLs(key string) []url.URL
```
URLs returns all URL values associated with the given key or returns the zero
value if the key does not exist or the value could not be parsed as a URL.

#### func (*Reader) URLsDefault

```go
func (this *Reader) URLsDefault(key string, Default []url.URL) []url.URL
```
URLsDefault returns all URL values associated with the given key or returns
provided defaults if the key does not exist or the values could not be parsed as
URLs.

#### func (*Reader) URLsError

```go
func (this *Reader) URLsError(key string) ([]url.URL, error)
```
URLsError returns all URL values associated with the given key with an error if
the key does not exist or the values could not be parsed as URLs.

#### func (*Reader) URLsFatal

```go
func (this *Reader) URLsFatal(key string) []url.URL
```
URLsFatal returns all URL values associated with the given key or calls
log.Fatal() if the key does not exist or the values could not be parsed as URLs.

#### func (*Reader) URLsPanic

```go
func (this *Reader) URLsPanic(key string) []url.URL
```
URLsPanic returns all URL values associated with the given key or panics if the
key does not exist or the values could not be parsed as URLs.

#### type Source

```go
type Source interface {
	Initialize()
	Strings(key string) ([]string, error)
}
```

Source defines the methods required by a Reader. The Strings method returns all
values associated with the given key with an error if the key does not exist.

#### func  FromRequiredInProductionJSONFile

```go
func FromRequiredInProductionJSONFile(filename string) Source
```
If we are running in a production (non-development) environment the specified
JSON filename is required; otherwise it's optional.