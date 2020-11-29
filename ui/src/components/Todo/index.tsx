import React from 'react';
import { useRouteMatch, withRouter } from 'react-router';

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
      selectable={true}
    >
      Todo: {url}
    </Heading>
  );
};

export default Todo;
