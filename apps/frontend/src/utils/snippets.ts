export const generateSnippet = (language: 'python' | 'curl' | 'node' | 'env', port: number, model: string = 'gemini-1.5-pro'): string => {
    const baseUrl = `http://localhost:${port}/v1`;
    const apiKey = 'dummy'; // Proxy handles auth to Google, but client might need a dummy key

    switch (language) {
        case 'python':
            return `from openai import OpenAI

client = OpenAI(
    base_url="${baseUrl}",
    api_key="${apiKey}"
)

response = client.chat.completions.create(
    model="${model}",
    messages=[{"role": "user", "content": "Hello, Siberia!"}]
)

print(response.choices[0].message.content)`;

        case 'node':
            return `import OpenAI from 'openai';

const client = new OpenAI({
    baseURL: '${baseUrl}',
    apiKey: '${apiKey}',
});

async function main() {
    const response = await client.chat.completions.create({
        model: '${model}',
        messages: [{ role: 'user', content: 'Hello, Siberia!' }],
    });
    console.log(response.choices[0].message.content);
}

main();`;

        case 'curl':
            return `curl ${baseUrl}/chat/completions \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer ${apiKey}" \\
  -d '{
    "model": "${model}",
    "messages": [
      {
        "role": "user",
        "content": "Hello, Siberia!"
      }
    ]
  }'`;

        case 'env':
            return `OPENAI_BASE_URL=${baseUrl}
OPENAI_API_KEY=${apiKey}`;

        default:
            return '';
    }
};
