const { makeExtendSchemaPlugin, gql } = require("graphile-utils");

const MyPlugin = makeExtendSchemaPlugin(build => {
  // Get any helpers we need from `build`
  const { pgSql: sql, inflection } = build;

  return {
    typeDefs: gql`
      extend type Query {
        hello: String!
      }

      input WebEventInput {
        stationId: Int!
        command: String!
        args: String
      }

      type WebEventPayload {
        result: String!
        eventId: Int!
      }

      extend type Mutation {
        postWebEvent(input: WebEventInput!): WebEventPayload
      }
    `,
    resolvers: {
      Query: {
        hello: async () => {
          console.log('hello')
          return 'world!';
        },
      },
      Mutation: {
        postWebEvent: async (_query, args, context, resolveInfo) => {
          const { pgClient } = context;

          const { input } = args;

          const data = {
            via: 'web', // redundant?
            from: 'anonymous',
            tokens: [input.command.trim(), input.args.trim()]
          }

          const {
            rows: [event]
          } = await pgClient.query(`
            insert into events (station_id, name, data)
              values ($1::integer, $2::text, $3::jsonb)
              returning *`, [input.stationId, 'WEB_COMMAND', data])

          return {
            result: 'OK',
            eventId: event.id,
          }
        }
      },
    },
  };
});

module.exports = MyPlugin; 
