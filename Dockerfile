# =========================
# Etapa 1 — Build
# =========================
FROM gradle:8.5-jdk21 AS builder

WORKDIR /app

# Copia apenas arquivos necessários para resolver dependências
COPY build.gradle settings.gradle gradlew ./
COPY gradle gradle
RUN chmod +x gradlew
# Baixa dependências (cache eficiente)
RUN ./gradlew dependencies --no-daemon

# Copia o restante do projeto
COPY src src

# Gera o JAR
RUN ./gradlew bootJar --no-daemon


# =========================
# Etapa 2 — Runtime
# =========================
FROM eclipse-temurin:21-jre-alpine

WORKDIR /app

# Copia o JAR gerado na etapa de build
COPY --from=builder /app/build/libs/*.jar app.jar

# Porta padrão usada pelo Render
EXPOSE 8080

# Render define a variável PORT automaticamente
ENV PORT=8080

# Comando de inicialização
ENTRYPOINT ["java", "-jar", "app.jar"]
