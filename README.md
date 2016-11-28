# You gotta keep them seperated

Credentials and state. It's inevitable when you use Stack Forming tools like CloudFormation, Bosh, and Terraform, that you will have some state you need to maintain. This state either wants you to provide credentials, or will helpfully create them for you. But, you don't want to store the credentials in Github, they share repo's by mistake, and so might you.

I use 2FA, specifically Lastpass. But given any source of configuration, which you want to manage, and credentials, which you want to store, how do you merge and split the two easily?


# Paths for injection: 

## 1. I have the credentials:

If I have the credentials, I want to inject the creds into a known place in a file. I add a moustache and wham!

Given:
 
  * Credentials as string
  * a file or stdin of the statefile with moustache {{placeholder}}
  * moustache name
  
Output:

  * finalized statefile

## 2. I do not have the credentials or state file 

If I do not have either the creds or the state file (or the statefile is somewhat empty), I still want to act like I have the statefile when I update a resource.

1. Input file does not match moustache, then just output the input file.
2. If the Input file does match the moustache, but no creds, error out.
3. If there is not input or input file, do not create output or output file


## Usage

```
offspring-state inject -k the-credentials -f test/cred -s test/state-file.thing

and 

cat test/state-file.thing | offspring-state inject -k the-credentials -f test/cred -s test/state-file.thing
```

**TODO**: Now we also want to try using <() etc.



# Extraction

Given:
  
  * statefile
  * extraction regex with matching group
  * moustache name
    
Output:

  * statefile with moustache name in place
  * credential as stdout or file
  * error if could not match, or match empty


