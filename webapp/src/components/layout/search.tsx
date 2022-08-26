import { useRouter } from "next/router";

import { SearchURL } from "#lib/urls";
import {
  FormClient,
  ISubmitEvent,
  IFormElements,
  FormSection,
} from "#components/forms";
import { Search } from "#components/inputs";

// eslint-disable-next-line
import styles from "./search.module.scss";
import { LinkInternal } from "#components/links";

export function GlobalSearch() {
  const router = useRouter();
  const formID = "global-search";

  async function handleSubmit(event: ISubmitEvent) {
    const query = (event.currentTarget.elements as IFormElements<"query">)[
      "query"
    ].value;
    router.push(new SearchURL(query).toString());
  }

  return (
    <FormClient id={formID} className={styles.block} onSubmit={handleSubmit}>
      <Search className={styles.search} id={`${formID}-query`} name="query">
        Search
      </Search>
      <FormSection className={styles.advanced}>
        <LinkInternal href={"/search"} target="_blank">
          Advanced
        </LinkInternal>
      </FormSection>
    </FormClient>
  );
}
