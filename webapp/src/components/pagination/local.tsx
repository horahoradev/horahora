import { useRef, useState } from "react";

import { FormClient, FormSection } from "#components/forms";
import {
  IListUnorderedProps,
  ListItem,
  ListUnordered,
} from "#components/lists";
import { blockComponent } from "#components/meta";
import { type IPagination } from "#lib/pagination";
import { NumberInput } from "#components/inputs";
import { Button, ButtonSubmit } from "#components/buttons";

// eslint-disable-next-line
import styles from "./internal.module.scss";

export interface IPaginationLocalProps extends IListUnorderedProps {
  pagination: Omit<IPagination, "limit">;
  onPageChange: (page: number) => Promise<void>;
}

export const PaginationLocal = blockComponent(styles.block, Component);

function Component({
  pagination,
  onPageChange,
  ...blockProps
}: IPaginationLocalProps) {
  const [isChangingPage, switchPageChange] = useState(false);
  const { totalCount } = pagination;
  const limit = 50;
  const totalPages = pagination.totalPages ?? Math.ceil(totalCount / limit);
  const currentPage = pagination.currentPage ?? totalPages;
  const isLastPage = currentPage === totalPages;
  const prevPage = currentPage - 1;
  const nextPage = currentPage + 1;

  async function changePage(page: number) {
    if (isChangingPage || page !== currentPage) {
      return;
    }

    const isValidPage = page > 0 && page <= totalPages!;

    if (!isValidPage) {
      throw new Error(`Invalid page of "${page}"`);
    }

    try {
      switchPageChange(true);
      await onPageChange(page);
    } finally {
      switchPageChange(false);
    }
  }

  return (
    <ListUnordered {...blockProps}>
      <ListItem>
        {currentPage === 1 ? (
          "..."
        ) : (
          <Button
            onClick={async () => {
              await changePage(1);
            }}
          >
            1
          </Button>
        )}
      </ListItem>

      <ListItem>
        {prevPage <= 1 ? (
          "..."
        ) : (
          <Button
            onClick={async () => {
              await changePage(prevPage);
            }}
          >
            {prevPage}
          </Button>
        )}
      </ListItem>

      <ListItem>
        <CurrentPage
          currentPage={currentPage}
          totalPages={totalPages}
          onPageChange={changePage}
        />
      </ListItem>

      <ListItem>
        {nextPage >= totalPages ? (
          "..."
        ) : (
          <Button
            onClick={async () => {
              await changePage(nextPage);
            }}
          >
            {nextPage}
          </Button>
        )}
      </ListItem>

      <ListItem>
        {isLastPage ? (
          "..."
        ) : (
          <Button
            onClick={async () => {
              await changePage(totalPages);
            }}
          >
            {totalPages}
          </Button>
        )}
      </ListItem>
    </ListUnordered>
  );
}

interface ICurrentPageProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => Promise<void>;
}

function CurrentPage({
  currentPage,
  totalPages,
  onPageChange,
}: ICurrentPageProps) {
  const inputRef = useRef<HTMLInputElement>(null);

  async function handleSubmit() {
    if (!inputRef.current) {
      return;
    }

    const selectedPage = Number(inputRef.current.value);

    if (selectedPage === currentPage) {
      return;
    }

    await onPageChange(selectedPage);
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
