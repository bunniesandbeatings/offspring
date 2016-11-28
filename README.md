# You gotta keep them separated

Credentials and state. It's inevitable when you use Stack Forming tools like CloudFormation, Bosh, and Terraform, 
that you will have some state you need to maintain. This state either wants you to provide credentials, or will
helpfully create them for you. But, you don't want to store the credentials in GitHub, they share repo's by mistake,
and so might you.

I use 2FA, specifically LastPass. But given any source of configuration, which you want to manage, and credentials,
which you want to store, how do you merge and split the two easily?

# All pipes

You can specify the input file if you must, but this tool is designed to:
  * pipe a state-file through a set of `inject` or `extract` commands.
  * use temprorary pipes `<()` `>()` to avoid putting credentials on disk

# Injection: 

## Path: I have the credentials

If I have the credentials, I want to inject the creds into a known place in a file. I add a moustache and wham!

Given:
 
  * Credentials as string
  * a file or stdin of the state-file with moustache {{placeholder}}
  * moustache name
  
Output:

  * finalized state-file content on stdout.
  * User uses redirection to commit state, or pipes to process further

## Path: I do not have the credentials or state file 

If I do not have either the creds or the state file (or the state-file is somewhat empty), I still want to act like I
have the state-file when I update a resource.

  * Input file does not match moustache, then just output the input file.

## Usage

```
offspring-state inject -k the-credentials -f test/cred -s test/state-file.thing > outfile.txt
```

and
 
```
cat test/state-file.thing | offspring-state inject -k the-credentials -f test/cred -s test/state-file.thing > outfile.txt
```

# Extraction

Given:
  
  * state-file
  * extraction regex with matching group
  * moustache name
    
Output:

  * state-file with moustache name in place
  * credential as stdout or file
  * error if could not match, or match empty

## Usage

```
offspring-state extract -k the-credential -f some-creds.file \
  -p $'(?s)match-me:.*?\'(-----BEGIN RSA PRIVATE KEY-----.*?-----END RSA PRIVATE KEY-----.*?)\'' \
  > outfile.txt
```

Again stdin also works.

## Regex Hints
Working with regex is a pain, but it's the only reliable way I could find for ensuring the accurate captures. 
Luckily the `offspring-state` tool is designed to chain, making testing your regex's a little simpler.

It's important to learn about setting flags in [golang regular expressions](https://golang.org/pkg/regexp/syntax/).
Most keys are multiline, so this pattern is your friend:

`(?s)match-me:.*?'(-----BEGIN RSA PRIVATE KEY-----.*?-----END RSA PRIVATE KEY-----.*?)'`

**Remember**:

  * Single-quotes around multiline credentials make YAML deal with them without needing indentation.
  * Single-quotes are a pain on the command line, but in bash `$''` is a great workaround.
  * Capture enough to be certain you got the credential you wanted.
  * Use `(?s)`, causing `.` to consume newlines.
  * Use `.*?` to negate greediness because of `(?s)`.

**Full Example**: 

```
cat test/state-file.thing | \
offspring-state inject -k the-credentials -f test/cred | \
offspring-state extract -k the-credential \
  -p $'(?s)match-me:.*?\'(-----BEGIN RSA PRIVATE KEY-----.*?-----END RSA PRIVATE KEY-----.*?)\''
```

# Testing:

  Currently the only test is install:
  
  ```
  go install .
  ```
  
  Ensure interpolation (two identical keys with newline before closing quote):

  ```
  I am a strange
  Kind of state file
  match-me: '-----BEGIN RSA PRIVATE KEY-----
  MIIEowIBAAKCAQEA1LLN4YbjcNE4cf9OpFERq+xUd3CAiIrzlAH7u/lLMoU2Ssko
  ... snip ...
  N5/jRa/s4Eq9FFxGnCPMy1tLcsifj4mJzxUMN/efNKvxH9BMdLjI
  -----END RSA PRIVATE KEY-----
  '
  no-match: '-----BEGIN RSA PRIVATE KEY-----
  MIIEowIBAAKCAQEA1LLN4YbjcNE4cf9OpFERq+xUd3CAiIrzlAH7u/lLMoU2Ssko
  ... snip ...
  N5/jRa/s4Eq9FFxGnCPMy1tLcsifj4mJzxUMN/efNKvxH9BMdLjI
  -----END RSA PRIVATE KEY-----
  '
  ```
  
  Ensure reversability (only one key file is replaced with the template variable)
  
  ```
  go install .; cat test/state-file.thing | \
    offspring-state inject -k the-credentials -f test/cred | \
    offspring-state extract -k the-credentials \
      -p $'(?sm)match-me:.*?\'(-----BEGIN RSA PRIVATE KEY-----.*?-----END RSA PRIVATE KEY-----.*?)\''
  I am a strange
  Kind of state file
  match-me: '{{the-credentials}}'
  no-match: '-----BEGIN RSA PRIVATE KEY-----
  MIIEowIBAAKCAQEA1LLN4YbjcNE4cf9OpFERq+xUd3CAiIrzlAH7u/lLMoU2Ssko
  ... snip ...
  N5/jRa/s4Eq9FFxGnCPMy1tLcsifj4mJzxUMN/efNKvxH9BMdLjI
  -----END RSA PRIVATE KEY-----
  '
  ```
  
  There should be NO diff:
  
  ```
  go install .; cat test/state-file.thing | \
    offspring-state inject -k the-credentials -f test/cred | \
    offspring-state extract -k the-credentials \
      -p $'(?sm)match-me:.*?\'(-----BEGIN RSA PRIVATE KEY-----.*?-----END RSA PRIVATE KEY-----.*?)\'' | \
      diff - test/state-file.thing
  ```
  
  Try messing up the regex, you should see the diff.
  
# TODOs
  
  * try using <() for credential input
  * try using >() for credential output


