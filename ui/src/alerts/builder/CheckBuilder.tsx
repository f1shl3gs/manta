// Libraries
import React from 'react'

// Components
import CheckMetaCard from './CheckMetaCard'
import CheckConditionsCard from './CheckConditionsCard'

const CheckBuilder: React.FC = () => {
  return (
    <div className={'alert-builder'} data-testid={'query-builder'}>
      <CheckMetaCard />
      <CheckConditionsCard />
    </div>
  )
}

export default CheckBuilder
