export const UserRank = {
  REGULAR: 0,
  TRUSTED: 1,
  ADMIN: 2,
} as const;

export type IUserRank = typeof UserRank[keyof typeof UserRank]