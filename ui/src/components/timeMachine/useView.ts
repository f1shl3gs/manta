import { useState } from 'react';
import constate from 'constate';

import { ViewProperties, XYViewProperties } from 'types';

interface State {
  view: ViewProperties
}

const [ViewProvider, useLineView] = constate(
  (initialState: State) => {
    const [view, setView] = useState<ViewProperties>(initialState.view);

    return {
      view,
      setView
    };
  },
  value => {
    const { view, setView } = value;
    const properties = view as XYViewProperties;

    const onSetXColumn = (x: string) => {
      setView({
        ...properties,
        xColumn: x
      });
    };

    const onSetYColumn = (y: string) => {
      setView({
        ...properties,
        yColumn: y
      });
    };

    return {
      xColumn: properties.xColumn || 'xc',
      yColumn: properties.yColumn || 'yc',
      onSetXColumn,
      onSetYColumn,
      numericColumns: [] as string[]
    };
  }
);

export {
  ViewProvider,
  useLineView
};