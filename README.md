# GPT-CLI

> GPT CLI tool implemented with Golang

## Usage

```
    go install github.com/dkyaorui/gpt-cli@latest
```

### Config token

#### Set your private token

```
gpt-cli config token
```

#### View token has been configured

```
gpt-cli config token -s
```

### Config Model

#### Set GPT model

```
gpt-cli config model
```

#### View model has been configured

```
gpt-cli config model -s
```

## Library dependencies

1. github.com/fzdwx/infinite
2. github.com/sashabaranov/go-openai
3. github.com/spf13/cobra
4. github.com/spf13/viper
