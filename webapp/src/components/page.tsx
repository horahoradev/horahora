import { useState, useEffect } from "react";
import type { ReactNode } from "react";

import { Header } from "#components/header";
import { getUserdata } from "#api/index";

interface IPageProps {
  children: ReactNode;
}

export function Page({ children }: IPageProps) {
  const [userData, setUserData] = useState(null);

  useEffect(() => {
    let ignore = false;

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
      <Header userData={userData} />
      {children}
    </>
  );
}
