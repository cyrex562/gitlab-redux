import { GlFilteredSearchToken } from '@gitlab/ui';
import { s__, __ } from '~/locale';
import {
  OPERATORS_AFTER_BEFORE,
  OPERATORS_IS,
} from '~/vue_shared/components/filtered_search_bar/constants';
import DateToken from '~/vue_shared/components/filtered_search_bar/tokens/date_token.vue';

// Token types
export const FEED_TOKEN = 'feedToken';
export const INCOMING_EMAIL_TOKEN = 'incomingEmailToken';
export const STATIC_OBJECT_TOKEN = 'staticObjectToken';

export const DEFAULT_SORT = { value: 'expires', isAsc: true };

export const SORT_OPTIONS = [
  {
    text: __('Name'),
    value: 'name',
    sort: {
      asc: 'name_asc',
      desc: 'name_desc',
    },
  },
  {
    text: __('Created date'),
    value: 'created',
    sort: {
      asc: 'created_asc',
      desc: 'created_desc',
    },
  },
  {
    text: __('Expiration date'),
    value: 'expires',
    sort: {
      asc: 'expires_at_asc_id_desc',
    },
  },
];

export const TOKENS = [
  {
    icon: 'key',
    title: s__('CredentialsInventory|Type'),
    type: 'filter',
    token: GlFilteredSearchToken,
    operators: OPERATORS_IS,
    unique: true,
    options: [
      {
        value: 'personal_access_tokens',
        title: s__('CredentialsInventory|Personal access tokens'),
      },
      { value: 'ssh_keys', title: s__('CredentialsInventory|SSH keys') },
      {
        value: 'resource_access_tokens',
        title: s__('CredentialsInventory|Project and group access tokens'),
      },
      { value: 'gpg_keys', title: s__('CredentialsInventory|GPG keys') },
    ],
  },
  {
    icon: 'status',
    title: s__('CredentialsInventory|State'),
    type: 'state',
    token: GlFilteredSearchToken,
    operators: OPERATORS_IS,
    unique: true,
    options: [
      { value: 'active', title: s__('CredentialsInventory|Active') },
      { value: 'inactive', title: s__('CredentialsInventory|Inactive') },
    ],
  },
  {
    icon: 'remove',
    title: s__('CredentialsInventory|Revoked'),
    type: 'revoked',
    token: GlFilteredSearchToken,
    operators: OPERATORS_IS,
    unique: true,
    options: [
      { value: 'true', title: __('Yes') },
      { value: 'false', title: __('No') },
    ],
  },
  {
    icon: 'history',
    title: s__('CredentialsInventory|Created date'),
    type: 'created',
    token: DateToken,
    operators: OPERATORS_AFTER_BEFORE,
    unique: true,
  },
  {
    icon: 'history',
    title: s__('CredentialsInventory|Expiration date'),
    type: 'expires',
    token: DateToken,
    operators: OPERATORS_AFTER_BEFORE,
    unique: true,
  },
  {
    icon: 'history',
    title: s__('CredentialsInventory|Last used date'),
    type: 'last_used',
    token: DateToken,
    operators: OPERATORS_AFTER_BEFORE,
    unique: true,
  },
];
