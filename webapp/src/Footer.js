import PrivacyPolicy from "./PrivacyPolicy";
import TermsOfService from "./TOS";

function Footer(props) {
  const { userData, dataless } = props;

  return (
    <nav className="h-8 w-full bg-white shadow">
        <div className="flex justify-around">
            <div><PrivacyPolicy></PrivacyPolicy></div>
            <div><TermsOfService></TermsOfService></div>
        </div>
    </nav>
  );
}

export default Footer;
