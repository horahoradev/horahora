import { useRouter } from "next/router";
import { useRef } from "react";

import { FormClient, FormSection } from "#components/forms";
import { LinkInternal } from "#components/links";
import {
  IListUnorderedProps,
  ListItem,
  ListUnordered,
} from "#components/lists";
import { blockComponent } from "#components/meta";
import { IURLBuilder, type IPagination } from "#lib/pagination";
import { NumberInput } from "#components/inputs";

// eslint-disable-next-line
import styles from "./internal.module.scss";
import { Button, ButtonSubmit } from "#components/buttons";

export interface IPaginationInternalProps extends IListUnorderedProps {
  pagination: IPagination;
  urlBuilder: IURLBuilder;
}

export const PaginationInternal = blockComponent(styles.block, Component);

function Component({
  pagination,
  urlBuilder,
  ...blockProps
}: IPaginationInternalProps) {
  const { limit, totalCount } = pagination;
  let { currentPage, totalPages } = pagination;

  if (!totalPages) {
    totalPages = Math.ceil(totalCount / limit);
  }

  if (!currentPage) {
    currentPage = totalPages;
  }

  const isLastPage = currentPage === totalPages;
  const prevPage = currentPage - 1;
  const nextPage = currentPage + 1;

  return (
    <ListUnordered {...blockProps}>
      <ListItem>
        {currentPage === 1 ? (
          "..."
        ) : (
          <LinkInternal href={urlBuilder(1)}>1</LinkInternal>
        )}
      </ListItem>

      <ListItem>
        {prevPage <= 1 ? (
          "..."
        ) : (
          <LinkInternal href={urlBuilder(prevPage)}>{prevPage}</LinkInternal>
        )}
      </ListItem>

      <ListItem>
        <CurrentPage
          currentPage={currentPage}
          totalPages={totalPages}
          urlBuilder={urlBuilder}
        />
      </ListItem>

      <ListItem>
        {nextPage >= totalPages ? (
          "..."
        ) : (
          <LinkInternal href={urlBuilder(nextPage)}>{nextPage}</LinkInternal>
        )}
      </ListItem>

      <ListItem>
        {isLastPage ? (
          "..."
        ) : (
          <LinkInternal href={urlBuilder(totalPages)}>
            {totalPages}
          </LinkInternal>
        )}
      </ListItem>
    </ListUnordered>
  );
}

interface ICurrentPageProps {
  currentPage: number;
  totalPages: number;
  urlBuilder: IURLBuilder;
}

function CurrentPage({
  currentPage,
  totalPages,
  urlBuilder,
}: ICurrentPageProps) {
  const router = useRouter();
  const inputRef = useRef<HTMLInputElement>(null);

  async function handleSubmit() {
    if (!inputRef.current) {
      return;
    }

    const selectedPage = Number(inputRef.current.value);

    if (selectedPage === currentPage) {
      return;
    }

    router.push(urlBuilder(selectedPage));
  }

  return (
    <FormClient id="pagination" onSubmit={handleSubmit} isSubmitSection={false}>
      <NumberInput
        id="pagination-page"
        name="page"
        min={1}
        max={totalPages}
        step={1}
        defaultValue={currentPage}
        inputRef={inputRef}
      />
      <FormSection>
        <Button
          onClick={() => {
            inputRef.current?.stepDown();
          }}
        >
          -1
        </Button>
      </FormSection>
      <FormSection>
        <Button
          onClick={() => {
            inputRef.current?.stepUp();
          }}
        >
          +1
        </Button>
      </FormSection>
      <FormSection>
        <ButtonSubmit>Go!</ButtonSubmit>
      </FormSection>
    </FormClient>
  );
}
