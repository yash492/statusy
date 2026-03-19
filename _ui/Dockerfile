FROM oven/bun:1.3-alpine AS builder

WORKDIR /app

COPY package.json ./
COPY bun.lock ./
COPY src ./

RUN bun install
COPY . .
RUN bun run build

FROM oven/bun:1.3-alpine
COPY --from=builder /app/build build/
COPY --from=builder /app/package.json .
EXPOSE 3000/tcp

# Run the SvelteKit application using Bun
CMD ["bun", "build/index.js"]

