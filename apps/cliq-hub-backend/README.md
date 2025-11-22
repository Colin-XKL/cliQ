## cliq-hub-backend

A new Golang backend that generates `.cliqfile.yaml` templates via LLM.

- Location: `cliq-hub-backend/`
- Language: Go 1.22
- Dependencies: `gin`, `go-openai`, `yaml.v3`

### Run

1. Set environment variables:

   - `LLM_API_KEY`: your OpenAI-compatible API key
   - `LLM_MODEL`: e.g. `gpt-4o-mini`
   - `LLM_BASE_URL` (optional): custom base URL for OpenAI-compatible endpoints

2. Start server:

```
LLM_API_KEY=sk-... LLM_MODEL=gpt-4o-mini go run ./cliq-hub-backend/cmd/server
```

Server listens on `:8080` by default (set `PORT` to override).

### API

- `POST /v1/templates/generate`

Request `application/json`:

```
{
  "command_example": "pngquant input.png --output output.png",
  "name": "PNGQuant 压缩工具",
  "description": "使用 pngquant 高效压缩 PNG 图片",
  "author": "user123",
  "encoding": "plain" // or "base64"
}
```

Response `application/json` embeds YAML content:

```
{
  "yaml": "name: PNGQuant 压缩工具\n...",
  "encoding": "plain",
  "meta": {
    "name": "PNGQuant 压缩工具",
    "version": "1.0",
    "cliq_template_version": "1.0"
  }
}
```

If `encoding` is `base64`, `yaml` contains the Base64-encoded YAML string.

Errors:

- `400` invalid input
- `502` LLM failure or unusable output
- `422` YAML parses but fails validation