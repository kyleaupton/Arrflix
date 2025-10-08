import path from 'node:path'
import url from 'node:url'
import { defineConfig } from '@hey-api/openapi-ts'

const dirname = path.dirname(url.fileURLToPath(import.meta.url))
const input = path.join(dirname, '..', 'backend/internal/http/docs/swagger.json')

export default defineConfig({
  input,
  output: 'src/client',
  plugins: ['@tanstack/vue-query'],
})
