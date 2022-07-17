import { useState, useEffect, type ReactNode } from "react";
import { useRouter } from "next/router";

import { Header } from "#components/header";
import { getUserdata } from "#api/index";

interface IPageProps {
  children: ReactNode;
}

export function Page({ children }: IPageProps) {
  const router = useRouter();
  const [userData, setUserData] = useState<Record<string, unknown>>();

  useEffect(() => {
    let ignore = false;

    if (router.pathname.startsWith("/authentication")) {
      return () => {
        ignore = true;
      };
    }

    (async () => {
      let userData = await getUserdata();
      if (!ignore) setUserData(userData);
    })();

    return () => {
      ignore = true;
    };
  }, []);

  return (
    <>
      <Header userData={userData} dataless={!userData} />
      {children}
    </>
  );
}
