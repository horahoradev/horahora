import { type ReactNode } from "react";

interface IPageProps {
  children: ReactNode;
}

export function Page({ children }: IPageProps) {
  return <>{children}</>;
}
