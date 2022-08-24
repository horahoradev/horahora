import { Page } from "#components/page";
import { LinkInternal } from "#components/links";
import { ListItem, ListUnordered } from "#components/lists";
import { Article } from "#components/article";

function ArchivesPage() {
  return (
    <Page title="Archives">
      <Article>
        <ListUnordered>
          <ListItem>
            <LinkInternal href="/account/archives/requests">
              Requests
            </LinkInternal>
          </ListItem>
          <ListItem>
            <LinkInternal href="/account/archives/events">Events</LinkInternal>
          </ListItem>
          <ListItem>
            <LinkInternal href="/account/archives/downloads">
              Download Progress
            </LinkInternal>
          </ListItem>
        </ListUnordered>
      </Article>
    </Page>
  );
}

export default ArchivesPage;
