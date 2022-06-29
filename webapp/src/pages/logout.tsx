import { useRouter } from "next/router";
import { useEffect } from "react";

import { postLogout } from "#api/index";

function LogoutPage() {
  const router = useRouter();

  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      await postLogout();
      if (!ignore) router.push("/");
    };

    fetchData();

    return () => {
      ignore = true;
    };
  }, [router]);

  return <></>;
}

export default LogoutPage;
