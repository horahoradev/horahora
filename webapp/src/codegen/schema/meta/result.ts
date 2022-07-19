/*
  This module was created by the codegen, do not edit it manually.
*/
/**
 * A local copy of `draft-07` meta schema.
 */
export type IMetaSchema = (IMetaSchema1 & IMetaSchema2)
export type NonNegativeInteger = number
export type NonNegativeIntegerDefault0 = NonNegativeInteger
export type IMetaSchema2 = ({
$id?: string
$schema?: string
$ref?: string
$comment?: string
title?: string
description?: string
default?: true
readOnly?: boolean
writeOnly?: boolean
examples?: true[]
multipleOf?: number
maximum?: number
exclusiveMaximum?: number
minimum?: number
exclusiveMinimum?: number
maxLength?: NonNegativeInteger
minLength?: NonNegativeIntegerDefault0
pattern?: string
additionalItems?: IMetaSchema2
items?: (IMetaSchema2 | SchemaArray)
maxItems?: NonNegativeInteger
minItems?: NonNegativeIntegerDefault0
uniqueItems?: boolean
contains?: IMetaSchema2
maxProperties?: NonNegativeInteger
minProperties?: NonNegativeIntegerDefault0
required?: StringArray
additionalProperties?: IMetaSchema2
definitions?: {
[k: string]: IMetaSchema2
}
properties?: {
[k: string]: IMetaSchema2
}
patternProperties?: {
[k: string]: IMetaSchema2
}
dependencies?: {
[k: string]: (IMetaSchema2 | StringArray)
}
propertyNames?: IMetaSchema2
const?: true
/**
 * @minItems 1
 */
enum?: [true, ...(unknown)[]]
type?: (SimpleTypes | [SimpleTypes, ...(SimpleTypes)[]])
format?: string
contentMediaType?: string
contentEncoding?: string
if?: IMetaSchema2
then?: IMetaSchema2
else?: IMetaSchema2
allOf?: SchemaArray
anyOf?: SchemaArray
oneOf?: SchemaArray
not?: IMetaSchema2
} | boolean)
/**
 * @minItems 1
 */
export type SchemaArray = [IMetaSchema2, ...(IMetaSchema2)[]]
export type StringArray = string[]
export type SimpleTypes = ("array" | "boolean" | "integer" | "null" | "number" | "object" | "string")

export interface IMetaSchema1 {
$id?: string
$schema?: string
$ref?: string
$comment?: string
title?: string
description?: string
default?: true
readOnly?: boolean
writeOnly?: boolean
examples?: true[]
multipleOf?: number
maximum?: number
exclusiveMaximum?: number
minimum?: number
exclusiveMinimum?: number
maxLength?: NonNegativeInteger
minLength?: NonNegativeIntegerDefault0
pattern?: string
additionalItems?: IMetaSchema2
items?: (IMetaSchema2 | SchemaArray)
maxItems?: NonNegativeInteger
minItems?: NonNegativeIntegerDefault0
uniqueItems?: boolean
contains?: IMetaSchema2
maxProperties?: NonNegativeInteger
minProperties?: NonNegativeIntegerDefault0
required?: StringArray
additionalProperties?: IMetaSchema2
definitions?: {
[k: string]: IMetaSchema2
}
properties?: {
[k: string]: IMetaSchema2
}
patternProperties?: {
[k: string]: IMetaSchema2
}
dependencies?: {
[k: string]: (IMetaSchema2 | StringArray)
}
propertyNames?: IMetaSchema2
const?: true
/**
 * @minItems 1
 */
enum?: [true, ...(unknown)[]]
type?: (SimpleTypes | [SimpleTypes, ...(SimpleTypes)[]])
format?: string
contentMediaType?: string
contentEncoding?: string
if?: IMetaSchema2
then?: IMetaSchema2
else?: IMetaSchema2
allOf?: SchemaArray
anyOf?: SchemaArray
oneOf?: SchemaArray
not?: IMetaSchema2
}


