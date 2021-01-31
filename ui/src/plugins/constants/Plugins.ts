// Markdown
import chronyMarkdown from 'plugins/markdown/chrony.md'

// Graphics
import chronyLogo from 'plugins/graphics/chrony.svg'

export interface PluginItem {
  id: string
  name: string
  url: string
  image?: string
  markdown?: string
}

const OTCL_PLUGINS_PATH = 'plugins'

export const OTCL_PLUGINS: PluginItem[] = [
  {
    id: 'chrony',
    name: 'Chrony',
    url: `${OTCL_PLUGINS_PATH}/chrony`,
    image: chronyLogo,
    markdown: chronyMarkdown,
  },
]
