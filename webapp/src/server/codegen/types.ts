export interface ICodegenModule extends Record<string, unknown> {
  default: ICodegenFunc;
}

export interface ICodegenFunc {
  (): Promise<ICodegen>
}

export interface ICodegen {
  /**
   * Export info used for creating index file.
   */
  exports: ICodegenExport;
  result: string
}

export interface ICodegenExport {
  types?: string[];
  concrete?: string[];
}
