/**
 * `onBlur` callback but only triggers the passed `callback`
 * when element's descendants aren't focused.
 * @template {Element} DOMInterface
 * @param {(...args: any[]) => void} callback A function to invoke on successful blur.
 * @param {any[]} [callbackArgs] Arguments to pass to the said function, if any.
 * @returns {(event: import("react").FocusEvent<DOMInterface>) => void}
 */
export function onParentBlur(
  callback,
  ...callbackArgs
) {
  return (event) => {
    
    if (event.currentTarget.contains(event.relatedTarget)) {
      event.preventDefault();
      return;
    }

    if (callbackArgs) {
      callback(...callbackArgs);
      return;
    }

    callback()
  };
}