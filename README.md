# CES-27 - Lab 01

Implementação do algoritmo de Ricart-Agrawala em Go.

Para fazer o *build* dos binários, rode `make all` no *root* do repositório. Ele irá criar uma pasta `target` com os dois binários do projeto.

O executável `sharedresource` não necessita de argumentos. O executável `process` necessita de um argumento com o número do processo e as portas de todos os processos. Exemplo para 3 processos (em terminais separados):

```
$ ./target/process 1 :10002 :10003 :10004
$ ./target/process 2 :10002 :10003 :10004
$ ./target/process 3 :10002 :10003 :10004
```

O número do processo precisa ser um número de 1 até o número de processos. A porta não pode ser `:10001` pois ela já é usada por `sharedresource`.