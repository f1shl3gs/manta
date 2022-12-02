// Libraries
import React, {useState} from 'react'

// Components
import ToolbarTab from 'src/timeMachine/PromqlToolbarTab'
import VariableToolbar from 'src/timeMachine/VariableToolbar'
import FunctionsToolbar from 'src/timeMachine/FunctionsToolbar'

type PromqlToolbarTabs = 'functions' | 'variables' | 'none'

const PromqlToolbar: React.FC = () => {
  const [activeTab, setActiveTab] = useState<PromqlToolbarTabs>('functions')
  const expanded = activeTab === 'functions' || activeTab === 'variables'

  let activeToolbar = activeTab === 'functions' ? <FunctionsToolbar /> : null
  if (activeTab === 'variables') {
    activeToolbar = <VariableToolbar />
  }

  return (
    <div className={'flux-toolbar'}>
      {expanded && (
        <div
          className={'flux-toolbar--tab-contents'}
          data-testid={`functions-toolbar-contents--${activeTab}`}
        >
          {activeToolbar}
        </div>
      )}

      <div className={'flux-toolbar--tabs'}>
        <ToolbarTab
          id={'functions'}
          name={'Functions'}
          active={activeTab === 'functions'}
          // @ts-ignore
          onClick={setActiveTab}
        />

        <ToolbarTab
          id={'variables'}
          name={'Variables'}
          // @ts-ignore
          onClick={setActiveTab}
          active={activeTab === 'variables'}
        />
      </div>
    </div>
  )
}

export default PromqlToolbar
