# You gotta keep them separated

Credentials and state. It's inevitable when you use Stack Forming tools like CloudFormation, Bosh, and Terraform, 
that you will have some state you need to maintain. This state either wants you to provide credentials, or will
helpfully create them for you. But, you don't want to store the credentials in GitHub, they share repo's by mistake,
and so might you.

I use 2FA, specifically LastPass. But given any source of configuration, which you want to manage, and credentials,
which you want to store, how do you merge and split the two easily?

# Paths for injection: 

## Path: I have the credentials

If I have the credentials, I want to inject the creds into a known place in a file. I add a moustache and wham!

Given:
 
  * Credentials as string
  * a file or stdin of the state-file with moustache {{placeholder}}
  * moustache name
  
Output:

  * finalized state-file

## Path: I do not have the credentials or state file 

If I do not have either the creds or the state file (or the state-file is somewhat empty), I still want to act like I
have the state-file when I update a resource.

1. Input file does not match moustache, then just output the input file.
2. If the Input file does match the moustache, but no creds, error out.
3. If there is not input or input file, do not create output or output file


## Usage

```
offspring-state inject -k the-credentials -f test/cred -s test/state-file.thing

and 

cat test/state-file.thing | offspring-state inject -k the-credentials -f test/cred -s test/state-file.thing
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


## Regex Hints
Working with regex is a pain, but the only good way I could come up with for ensuring the right captures. 
Luckily `offspring-state` tool is designed to chain, making testing your regex's a little simpler.

It's important to learn about setting flags in [golang regular expressions](https://golang.org/pkg/regexp/syntax/).
Most keys are multiline, so this pattern is your friend:

```
cat test/state-file.thing | \
offspring-state inject -k the-credentials -f test/cred | \
offspring-state extract -k the-credential \
  -p $'(?sm)match-me:.*?\'(-----BEGIN RSA PRIVATE KEY-----.*?-----END RSA PRIVATE KEY-----.*?)\''
```

Also note the `?` greedy negation

Lastly, dealing with single quoted strings for your regex:

```
echo $'It\'s Shell Programming'  # ksh, bash, and zsh only, does not expand variables
echo "It's Shell Programming"   # all shells, expands variables
echo 'It'\''s Shell Programming' # all shells, single quote is outside the quotes
echo 'It'"'"'s Shell Programming' # all shells, single quote is inside double quotes
```

# TODOs

  * biggest bug:
  
  ```
  go install .; cat test/state-file.thing | offspring-state inject -k the-credentials -f test/cred | offspring-state extract -k the-credentials -p $'(?sm)match-me:.*?\'(-----BEGIN RSA PRIVATE KEY-----.*?-----END RSA PRIVATE KEY-----.*?)\''
  ```
  
  currently replaces both creds because they look the same. I need an RE method that lets me replace a capture group.

  * try using <() for credential input
  * try using >() for credential output
  * support overriding input state-file
  * support writing to an alternate state-file




