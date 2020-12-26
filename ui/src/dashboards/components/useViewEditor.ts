import constate from 'constate';
import { useState } from 'react';

const [ViewEditorProvider, useViewEditor] = constate(
  () => {
    const [isViewingVisOptions, setIsViewingVisOptions] = useState(false);
    const onToggleVisOptions = () => setIsViewingVisOptions(!isViewingVisOptions)

    return {
      isViewingVisOptions,
      onToggleVisOptions
    };
  },
  value => value
);

export {
  ViewEditorProvider,
  useViewEditor
};