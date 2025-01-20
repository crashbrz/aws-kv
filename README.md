![License](https://img.shields.io/badge/license-sushiware-red)
![Issues open](https://img.shields.io/github/issues/crashbrz/aws-kv)
![GitHub pull requests](https://img.shields.io/github/issues-pr-raw/crashbrz/aws-kv)
![GitHub closed issues](https://img.shields.io/github/issues-closed-raw/crashbrz/aws-kv)
![GitHub last commit](https://img.shields.io/github/last-commit/crashbrz/aws-kv)

# AWS Key Validator

**aws-kv** is a tool for validating AWS credentials. It supports validating single credentials or bulk credentials from a file, with features for concurrency and debug mode. The tool leverages the AWS SDK to check the validity of credentials and provides detailed information about valid ones.

---

## Features

- Validate single AWS credentials (`AWS_KEY:secret`) via the `-k` flag.
- Validate multiple AWS credentials from a file via the `-f` flag.
- Concurrent processing with a configurable number of goroutines using the `-t` flag.
- Debug mode (`-d` flag) to display invalid credentials.
- Counters for the number of valid and invalid credentials.
- Displays detailed information about valid credentials, including Caller Identity and IAM details.

---

## Usage

### Command-Line Flags

| Flag      | Description                                                                                   |
|-----------|-----------------------------------------------------------------------------------------------|
| `-k`      | Provide a single AWS credential in the format `AWS_KEY:secret` for validation.                |
| `-f`      | Provide a file containing AWS credentials, one per line, in the format `AWS_KEY:secret`.      |
| `-t`      | Specify the number of goroutines for concurrent processing (default: 1).                      |
| `-d`      | Enable debug mode to display invalid credentials and count them in the summary.               |

### Examples

#### Validate a Single Credential
```bash
aws-kv -k "AKIAEXAMPLE:xyzSECRET" -t 5
```

#### Validate Credentials from a File
```bash
aws-kv -f credentials.txt -t 10
```

#### Enable Debug Mode
```bash
aws-kv -f credentials.txt -t 10 -d
```

---

## Output

### Valid Credentials

Valid credentials are displayed in **green**, along with detailed information:
```plaintext
Valid: AKIAEXAMPLE
Details:
Caller Identity: AIDAIEXAMPLE, Account: 123456789012, ARN: arn:aws:iam::123456789012:user/example

Number of valid credentials: 1
```

### Debug Mode (Invalid Credentials)
Invalid credentials are displayed in **red**:
```plaintext
Valid: AKIAEXAMPLE
Details:
Caller Identity: AIDAIEXAMPLE, Account: 123456789012, ARN: arn:aws:iam::123456789012:user/example

Invalid: AKIAINVALID

Number of valid credentials: 1
Number of invalid credentials: 1
```

---

## File Format for `-f`

Each line should contain one AWS credential in the format:
```plaintext
AWS_KEY:secret
```

Example:
```plaintext
AKIAEXAMPLE1:xyzSECRET1
AKIAEXAMPLE2:xyzSECRET2
AKIAEXAMPLE3:xyzSECRET3
```

---

## Building the Program
```bash
go build -o /openai-kv
```

---

## Cloning the Repository
```bash
git clone <repository-url>
cd <repository-name>
```

---

### License

aws-kv is licensed under the SushiWare license. For more information, check [docs/license.txt](docs/license.txt).
