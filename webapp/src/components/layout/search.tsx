import { useRouter } from "next/router";
import { faSearch } from "@fortawesome/free-solid-svg-icons";

import { SearchURL } from "#lib/urls";
import {
  FormClient,
  ISubmitEvent,
  IFormElements,
  FormSection,
} from "#components/forms";
import { Search } from "#components/inputs";
import { LinkInternal } from "#components/links";
import { ButtonSubmit } from "#components/buttons";
import { Icon } from "#components/icons";

// eslint-disable-next-line
import styles from "./search.module.scss";

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
    <FormClient
      id={formID}
      className={styles.block}
      onSubmit={handleSubmit}
      isSubmitSection={false}
    >
      <Search className={styles.search} id={`${formID}-query`} name="query">
        Search
      </Search>
      <FormSection className={styles.advanced}>
        <LinkInternal href={"/search"} target="_blank">
          Advanced
        </LinkInternal>
      </FormSection>
      <FormSection>
        <ButtonSubmit className={styles.submit}>
          <Icon icon={faSearch} />
        </ButtonSubmit>
      </FormSection>
    </FormClient>
  );
}
