-- Autor: Joan Martorell Coll

with Ada.Text_IO; use Ada.Text_IO;
with Ada.Numerics.Discrete_Random;
with Random_Generators;

procedure main is

   package PontAuxiliar is
      protected type Pont is
         procedure Bloqueja(id: Integer; direccio: Character);
         entry CotxeAlNord(id: Integer; direccio: Character);
         entry CotxeAlSud(id: Integer; direccio: Character);
         procedure CotxeFora(id: Integer; direccio: Character);
         entry AmbulanciaDins(id: Integer);
         procedure AmbulanciaFora(id: Integer);
      private
         EsperaNord: Integer := 0;  -- Comptador de cotxes al nord
         EsperaSud: Integer := 0;   -- Comptador de cotxes al sud
         EnCurs: Boolean := False;   -- Vehicle passant pel pont
         Ambulancia: Boolean := False;   -- Ambulancia esperant
      end Pont;
   end PontAuxiliar;

   package body PontAuxiliar is

      protected body Pont is

         -- Vehicle arriba al pont
         procedure Bloqueja(id: Integer; direccio: Character) is
         begin
            -- Si es ambulancia...
            if id = 112 then
               Put_Line("+++++Ambulància " & Integer'Image(id) & " espera per entrar");
               Ambulancia := True;

            -- Altre cas
            else
               if direccio = 'S' then
                  EsperaSud := EsperaSud + 1;
                  Put_Line("El cotxe " & Integer'Image(id) & " espera a l'entrada SUD, esperen al SUD: " & Integer'Image(EsperaSud));
               else
                  EsperaNord := EsperaNord + 1;
                  Put_Line("El cotxe " & Integer'Image(id) & " espera a l'entrada NORD, esperen al NORD: " & Integer'Image(EsperaNord));
               end if;
            end if;
         end Bloqueja;

         -- Cotxes del nord esperant
         entry CotxeAlNord(id: Integer; direccio: Character) when not EnCurs and not Ambulancia and EsperaNord >= EsperaSud is
         begin
            EnCurs := True;
            EsperaNord := EsperaNord - 1;
            Put_Line("El cotxe " & Integer'Image(id) & " entra al pont. Esperen al NORD " & Integer'Image(EsperaNord));
         end CotxeAlNord;

         -- Cotxes del sud esperant
         entry CotxeAlSud(id: Integer; direccio: Character) when not EnCurs and not Ambulancia and EsperaSud > EsperaNord is
         begin
            EnCurs := True;
            EsperaSud := EsperaSud - 1;
            Put_Line("El cotxe " & Integer'Image(id) & " entra al pont. Esperen al SUD " & Integer'Image(EsperaSud));
         end CotxeAlSud;

         -- Cotxe surt del pont
         procedure CotxeFora(id: Integer; direccio: Character) is
         begin
            Put_Line("------> El vehicle " & Integer'Image(id) & " surt del pont");
            EnCurs := False;
         end CotxeFora;

         -- Aambulància passa pel pont
         entry AmbulanciaDins(id: Integer) when not EnCurs is
         begin
            EnCurs := True;
            Ambulancia := False;
            Put_Line("+++++Ambulància " & Integer'Image(id) & " és al pont");
         end AmbulanciaDins;

         -- Ambulància surt del pont
         procedure AmbulanciaFora(id: Integer) is
         begin
            Put_Line("------> El vehicle " & Integer'Image(id) & " surt del pont");
            EnCurs := False;
         end AmbulanciaFora;

      end Pont;

   end PontAuxiliar;

   Pont: PontAuxiliar.Pont;

   -- Task Vehicle
   task type Vehicle(id: Integer; direccio: Character);

   task body Vehicle is
   begin

      -- Si és ambulància...
      if id = 112 then

         delay Duration(Random_Generators.GeneradorNombreAleatori);

         Put_Line("L'ambulància " & Integer'Image(id) & " està en ruta");

         delay Duration(Random_Generators.GeneradorNombreGranAleatori);

         Pont.Bloqueja(id, direccio);
         Pont.AmbulanciaDins(id);

         delay Duration(Random_Generators.GeneradorNombreGranAleatori);

         Pont.AmbulanciaFora(id);

      -- Altre cas...
      else

         delay Duration(Random_Generators.GeneradorNombreAleatori);

         Put_Line("El cotxe " & Integer'Image(id) & " està en ruta en direcció " & Character'Image(direccio));

         delay Duration(Random_Generators.GeneradorNombreAleatori);

         Pont.Bloqueja(id, direccio);

         if direccio = 'N' then
            Pont.CotxeAlNord(id, direccio);
         else
            Pont.CotxeAlSud(id, direccio);
         end if;

         delay Duration(Random_Generators.GeneradorNombreGranAleatori);

         Pont.CotxeFora(id, direccio);
      end if;

   end Vehicle;

   -- Inicialització dels "processos"
   Cotxe1: Vehicle(1, 'S');
   Cotxe2: Vehicle(2, 'N');
   Cotxe3: Vehicle(3, 'S');
   Cotxe4: Vehicle(4, 'N');
   Cotxe5: Vehicle(5, 'S');
   Ambulancia: Vehicle(112, 'N');

begin

   null;

end main;
