# Sistema de Gerenciamento de Notas - Chaincode Hyperledger Fabric

Este projeto implementa um **chaincode** em Go para o gerenciamento de notas de alunos em disciplinas, utilizando o Hyperledger Fabric.

## Funcionalidades

- **Criar Nota**: Adiciona uma nova nota para um aluno em uma disciplina.
- **Consultar Nota**: Recupera os dados de uma nota específica.
- **Atualizar Nota**: Altera os dados de uma nota existente.
- **Deletar Nota**: Remove uma nota do ledger.
- **Verificar Existência**: Checa se uma nota existe pelo ID.

## Estrutura da Nota
A estrutura da nota é composta pelos seguintes campos:

- **ID**: Identificador único da nota.
- **AlunoID**: Identificador do aluno.
- **Disciplina**: Nome da disciplina.
- **Nota**: Valor da nota (float).
- **Data**: Data em que a nota foi registrada.
- **Timestamp**: Momento exato do registro ou atualização da nota.

## Exemplo de objeto Nota em JSON:
```json
   {
     "ID": "nota1",
     "AlunoID": "aluno123",
     "Disciplina": "Matemática",
     "Nota": 8.5,
     "Data": "2024-06-01",
     "Timestamp": "2024-06-01T14:30:00Z"
   }
   ```

## Como usar
1. Clone o repositório.
2. Instale as dependências com `go mod tidy`.
3. Importe o chaincode para sua rede Hyperledger Fabric.
4. Use os comandos apropriados para invocar as funções do chaincode.
