# Codegen

## Table of contents
- [Introduction](#introduction)
- [The rules of codegen](#the-rules-of-codegen)
- [Basic structure](#basic-structure)
    - [General](#general)
    - [Language-specific](#language-specific)
- [Codegen sections](#codegen-sections)
    - [JSON Schema](#json-schema)

## Introduction

The need of input validation requires using json schemas. But json schemas allow for more than that. Since they describe the shape of serializable data which can be expressed by interfaces, the interfaces can be derived from them too (which can be implemented as classes). All of that can be done manually as the need arises for sure, but manual approach is prone to errors and gets hard to keep track of and synchronize changes due to indirect relation between json schema files and the code using them. Therefore it is better to have a system which derives the needed code from the schemas.

## The rules of codegen
- **no runtime codegen**

    Because generated code can be quite verbose and complicated, it is pretty hard to debug the code which isn't expressed in the source files. And runtime code can't be statically analyzed consistently either so it's also a security issue.

- **no production codegen**

    Because production environment deals with real data and the code generated will inevitably manipulate it in some way with no easy way to revert, there must be not production-only codegen.

- **no automatic codegen**

    The developer has to initiate a command to update the generated code. If there is a system which can detect the desync, it has to log a warning along with the instructions on how to update.

- **all generated code is checked into the repo**

    This one allows to automatically check for consistency. Given the same input the codegen has to produce the same output and any difference is human error (i.e. the person comitting the code edited the output) or malicious intent.

Because almost all these points deal with security one way or another, it might give an impression that codegen is some sort of security hazard. But it's only due to its nature of creating arbitrary code in a server language. It is no more insecure than adding a 3rd-party dependency.

## Basic structure

### General
The codegen consist of 2 parts:

- codegen folder
- codegen lib

#### Codegen Folder

This is the folder where all generated code resides. This folder is further separated into codegen modules of this structure:

- generator file

    A file which exports a generator function with no arguments which returns a string of the code generated.

- result file

    A file into which generated string is written into

- index file

    If the language supports re-exports, the generated symbols are re-exported there.

For the folder to count as a codegen module, the generator file must be present. Codegen modules cannot have nested codegen modules within them.

#### Codegen Lib

This library is responsible for analyzing, formatting and writing the generated code. Its module can export whatever number of symbols to help wtih code generation, but has to export a function with this signature:
```typescript
async function runCodegen(): string[]
```

This function walks the `codegen folder` and performs these operations in a `codegen module`:
1. Runs the generator.
2. The resulting string then formatted and written to the `result file`.
3. `index file` is created and populated.
4. The folder path is added to the list of paths.

### Language-specific

#### Javascript
- folder: `/webapp/src/codegen`
- module:
    - index file: `_index.ts`
    - generator: `generator.ts`
    - result file: `result.ts`
- command: `npm run codegen`

#### Golang
- folder: `/cli/src/codegen`

    For now it lives within the `cli` until the structure for shared packages is figured out.

- module
    - generator: `generator.go`
    - result file: `exports.go`

    Golang doesn't seem to support re-exports so the result file is also an index file.

- command: `./horahora codegen`

## Codegen sections

These are sections all target language codegens have to implement.

### JSON Schema

- Schema folder - `/schema`
- Schema file ending - `.schema.json`

JSON schema codegen has to travel in the `schema folder` and for each `schema file` export these things:
- an interface of the target JSON
- an immutable variable storing the parsed version of schema
- a function to validate the target JSON

As an example Typescript generated exports would look like this:
```typescript
export interface IMeta {...};
export const metaSchema = {...} as const;
export function validateMeta(inputJSON: unknown): inputJSON is IMeta;
```
