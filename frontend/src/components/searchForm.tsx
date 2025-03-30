import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Form, FormControl, FormField, FormItem } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useNavigate } from "@tanstack/react-router";

const FormSchema = z.object({
  summonerName: z.string(),
});

export function SearchForm() {
  const navigate = useNavigate();

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      summonerName: "",
    },
  });

  function onSubmit({ summonerName }: z.infer<typeof FormSchema>) {
    const [name, tag] = summonerName.split("#");
    navigate({ to: "/summoner/$name/$tag", params: { name, tag } });
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="w-[300px]">
        <FormField
          control={form.control}
          name="summonerName"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <Input
                  autoComplete="off"
                  placeholder="Player name..."
                  {...field}
                />
              </FormControl>
            </FormItem>
          )}
        />
      </form>
    </Form>
  );
}
