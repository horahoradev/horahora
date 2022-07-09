import styles from "./footer.module.scss";

import { LinkInternal } from "#components/links";

interface IFooterProps extends Record<string, unknown> {}

export function Footer({ userData, dataless }: IFooterProps) {
  return (
    <footer className={styles.block}>
      <LinkInternal className={styles.link} href="/privacy-policy">Privacy Policy</LinkInternal>
      <LinkInternal className={styles.link} href="/terms-of-service">Terms of Service</LinkInternal>
    </footer>
  );
}
