-- Autor: Joan Martorell Coll

with Ada.Text_IO;
with Ada.Numerics.Discrete_Random;

package body Random_Generators is

   function GeneradorNombreAleatori return Integer is
      type RangAleatori is new Integer range 1..10;
      package Aleatori_Int is new Ada.Numerics.Discrete_Random(RangAleatori);
      use Aleatori_Int;
      Gen : Generator;
      Num : RangAleatori;
   begin
      Reset(Gen);
      Num := Random(Gen);
      return Integer(Num);
   end GeneradorNombreAleatori;

   function GeneradorNombreGranAleatori return Integer is
      type RangAleatori is new Integer range 4..6;
      package Aleatori_Int is new Ada.Numerics.Discrete_Random(RangAleatori);
      use Aleatori_Int;
      Gen : Generator;
      Num : RangAleatori;
   begin
      Reset(Gen);
      Num := Random(Gen);
      return Integer(Num);
   end GeneradorNombreGranAleatori;

end Random_Generators;
