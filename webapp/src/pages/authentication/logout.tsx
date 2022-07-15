import { useRouter } from "next/router";
import { useEffect } from "react";

import { logoutAccount } from "#api/authentication";

function LogoutPage() {
  const router = useRouter();

  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      await logoutAccount();
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
