import React from 'react';
import { useRouteMatch } from 'react-router-dom';

import {
  FontWeight,
  Heading,
  HeadingElement,
  Typeface,
} from '@influxdata/clockface';

const Todo: React.FC = () => {
  const { url } = useRouteMatch();

  return (
    <Heading
      element={HeadingElement.H1}
      type={Typeface.Rubik}
      weight={FontWeight.Bold}
      underline={false}
      selectable
    >
      Todo: {url}
    </Heading>
  );
};

export default Todo;
