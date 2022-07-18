# Codegen

## Table of contents
- [Introduction](#introduction)
- [The rules of codegen](#the-rules-of-codegen)

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

The codegen consist 3 parts:

- input files
- generator modules
- codegen lib
