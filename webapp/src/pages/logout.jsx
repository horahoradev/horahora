import { useRouter } from "next/router";
import { useEffect } from "react";

import * as API from "../api";

function LogoutPage() {
  const router = useRouter();

  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      await API.postLogout();
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
