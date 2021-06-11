import { useEffect } from "react";
import { useHistory } from "react-router-dom";

import { postLogout } from "./api";

function LogoutPage() {
  let history = useHistory();

  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      await postLogout();
      if (!ignore) history.push("/");
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [history]);

  return <></>;
}

export default LogoutPage;
