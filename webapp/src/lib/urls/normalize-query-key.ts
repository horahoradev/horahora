const NORMALIZATION_STRATEGY = {
  FIRST: "first",
  LAST: "last",
  ALL: "all",
} as const;

type IStrategy =
  typeof NORMALIZATION_STRATEGY[keyof typeof NORMALIZATION_STRATEGY];

interface INormalizeQueryKeyOptions {
  defaultValue?: string;
  strategy?: IStrategy;
  separator?: string;
}

const defaultNormalizeQueryKeyOptions: INormalizeQueryKeyOptions = {
  defaultValue: undefined,
  strategy: NORMALIZATION_STRATEGY.FIRST,
  separator: ",",
} as const;

export function normalizeQueryKey(
  query: string | string[] | undefined,
  options?: INormalizeQueryKeyOptions
) {
  const finalOptions = options
    ? { ...defaultNormalizeQueryKeyOptions, ...options }
    : defaultNormalizeQueryKeyOptions;
  const { defaultValue, strategy, separator } = finalOptions;
  let finalValue: typeof defaultValue;

  if (Array.isArray(query)) {
    switch (strategy) {
      case NORMALIZATION_STRATEGY.FIRST: {
        finalValue = query[0];
        break;
      }

      case NORMALIZATION_STRATEGY.LAST: {
        finalValue = query[query.length - 1];
        break;
      }

      case NORMALIZATION_STRATEGY.ALL: {
        finalValue = query.join(separator);
        break;
      }

      default: {
        throw new Error(`Unknown query normalization strategy "${strategy}"`);
      }
    }
  } else {
    finalValue = query;
  }

  if (!finalValue) {
    finalValue = defaultValue;
  }

  return finalValue;
}
