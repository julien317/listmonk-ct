<template>
  <section class="segments">
    <header class="columns page-header">
      <div class="column is-8">
        <h1 class="title is-4">
          Segments
          <span v-if="segments.length" class="has-text-grey-light">({{ segments.length }})</span>
        </h1>
        <p class="has-text-grey is-size-7">
          Saved, named subscriber queries you can reuse when sending a campaign to a segment of a list.
        </p>
      </div>
      <div class="column is-4 has-text-right">
        <b-button type="is-primary" icon-left="plus" @click="onNew">New segment</b-button>
      </div>
    </header>

    <b-table :data="segments" :loading="loading" hoverable>
      <b-table-column field="name" label="Name" v-slot="props">
        <a href="#" @click.prevent="onEdit(props.index)">{{ props.row.name }}</a>
      </b-table-column>
      <b-table-column field="query" label="Query" v-slot="props">
        <code class="is-size-7">{{ props.row.query }}</code>
      </b-table-column>
      <b-table-column label="" width="120" cell-class="has-text-right" v-slot="props">
        <a href="#" @click.prevent="onEdit(props.index)" aria-label="edit">
          <b-icon icon="pencil-outline" size="is-small" />
        </a>
        <a href="#" class="ml-4" @click.prevent="onDelete(props.index)" aria-label="delete">
          <b-icon icon="trash-can-outline" size="is-small" />
        </a>
      </b-table-column>

      <template #empty>
        <div class="has-text-grey has-text-centered p-6">
          No segments yet. Create one to reuse across campaigns.
        </div>
      </template>
    </b-table>

    <b-modal :active.sync="isModalOpen" :width="640" scroll="keep">
      <div class="modal-card" style="width: auto">
        <header class="modal-card-head">
          <p class="modal-card-title">{{ form.index === null ? 'New segment' : 'Edit segment' }}</p>
        </header>
        <section class="modal-card-body">
          <b-field label="Name" label-position="on-border">
            <b-input v-model="form.name" placeholder="e.g. French speakers" maxlength="200" required />
          </b-field>
          <b-field label="Query" label-position="on-border"
            message="Partial SQL over subscriber attributes. Example: subscribers.attribs->'tags' ? 'French'">
            <b-input v-model="form.query" type="textarea" rows="3"
              placeholder="subscribers.attribs->'tags' ? 'French'" />
          </b-field>
          <div class="mt-2">
            <b-button size="is-small" icon-left="magnify" :loading="testing" @click="onTest">Test</b-button>
            <span v-if="testCount !== null" class="ml-3 is-size-7 has-text-grey">
              {{ testCount }} subscribers match
            </span>
          </div>
        </section>
        <footer class="modal-card-foot has-text-right">
          <b-button @click="isModalOpen = false">Close</b-button>
          <b-button type="is-primary" :disabled="!form.name || !form.query" @click="onSaveModal">Save</b-button>
        </footer>
      </div>
    </b-modal>
  </section>
</template>

<script>
import Vue from 'vue';

export default Vue.extend({
  data() {
    return {
      segments: [],
      loading: false,
      isModalOpen: false,
      testing: false,
      testCount: null,
      form: { index: null, name: '', query: '' },
    };
  },

  methods: {
    loadSegments() {
      this.loading = true;
      this.$api.getSegments().then((data) => {
        this.segments = Array.isArray(data) ? data : [];
        this.loading = false;
      }).catch(() => {
        this.loading = false;
      });
    },

    onNew() {
      this.form = { index: null, name: '', query: '' };
      this.testCount = null;
      this.isModalOpen = true;
    },

    onEdit(index) {
      const s = this.segments[index];
      this.form = { index, name: s.name, query: s.query };
      this.testCount = null;
      this.isModalOpen = true;
    },

    onSaveModal() {
      const item = { name: this.form.name.trim(), query: this.form.query.trim() };
      if (this.form.index === null) {
        this.segments.push(item);
      } else {
        this.$set(this.segments, this.form.index, item);
      }
      this.persist();
      this.isModalOpen = false;
    },

    onDelete(index) {
      this.$utils.confirm('Delete this segment?', () => {
        this.segments.splice(index, 1);
        this.persist();
      });
    },

    persist() {
      this.$api.updateSegments(this.segments).then(() => {
        this.$utils.toast('Segments saved');
      });
    },

    onTest() {
      const q = this.form.query.trim();
      if (!q) {
        return;
      }
      this.testing = true;
      this.testCount = null;
      this.$api.getSubscribers({ query: q, per_page: 1 }).then((data) => {
        this.testCount = data.total;
        this.testing = false;
      }).catch(() => {
        this.testing = false;
      });
    },
  },

  mounted() {
    this.loadSegments();
  },
});
</script>
