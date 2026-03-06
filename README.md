# Proyecto 1: Expresiones Regulares y Autómatas
## Autor: José Antonio Mérida Castejón
### Descripción
Este proyecto consiste en la implementación de algoritmos básicos para la construcción de autómatas finitos a partir de expresiones regulares. Tiene como entrada expresiones regulares _r_ y cadenas _w_. A partir de cada expresión regular se transformá de su expresión infix a postfix, luego construyé un AFN utilizando construcción de Thompson. Este AFN se transforma a un AFD por construcción de subconjuntos, y por último se minimiza la cantidad de estados. Con cada uno de estos autómatas se determina si la cadena _w_ pertenece a _L(r)_
### Ubicación de los archivos
- **DATA**: El archivo .txt de expresiones regulares que funcionan como entrada
- **OUT**: Folders conteniendo un PDF ilustrando NFA, DFA y DFA con estados minimizados para cada expresión regular dentro del archivo
### Dependencias
Para correr este programa se debe de tener instalado [Graphviz](https://graphviz.org/) para el renderizado de los autómatas y [Golang](https://go.dev/doc/install).

### Demostración
Archivo de entrada regex.txt
```
(a|b)*abb(a|b)*
b+abc+
```

Salida en Consola / Conversion a Postfix

![image](https://github.com/user-attachments/assets/3d39a40f-b454-411e-bc2b-1a72b5a27134)

DFA Minimizado (a|b)\*abb(a|b)\*

![image](https://github.com/user-attachments/assets/5812c2ac-e668-42e8-a301-2dbe1f6a54d4)

DFA Minimizado b+abc+

![image](https://github.com/user-attachments/assets/270f2319-2b48-4b0d-8992-c184a935236b)


### Corrida / Compilacion
```
go run cmd/main.go
```

```
go build -o bin/main.exe cmd/main.go
```

```
./bin/main
```

### Ejemplificación de Manejo de Errores / Validación de Regex (Obligado por Bidkar :sob:)
Balanceo de Expresion
```
)aa(b)
```
![image](https://github.com/user-attachments/assets/b236a419-d948-47d7-bfb1-5804f9d296ac)

Balanceo de Expresion con Caracteres Escapados
```
\)aa(b)
\(F)
```
![image](https://github.com/user-attachments/assets/162a8947-4ab4-4869-ae37-43d28038c996)

Expresiones Regulares No Validas (Error en Construcción con Thompson)
```
(a|)
```
![image](https://github.com/user-attachments/assets/322ecd9f-b506-4160-bdd8-b29007acf725)
```
***
```
![image](https://github.com/user-attachments/assets/b0e28d1b-d075-4c9e-9d3a-8097b7168145)
